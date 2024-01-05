use anyhow::{Context, Result};
use axum::extract::State;
use axum::http::{HeaderMap, StatusCode};
use axum::routing::post;
use axum::{Json, Router};
use regex::Regex;
use scryfall::card::Card;
use serde::{Deserialize, Serialize};
use std::env;
use std::fmt::Display;
use std::sync::Arc;
use teloxide::payloads::SendMessage;
use teloxide::prelude::*;
use teloxide::requests::JsonRequest;
use teloxide::types::ParseMode;
use teloxide::utils::markdown::{escape, escape_link_url};
use teloxide::Bot;
use tokio::net::TcpListener;

#[derive(Clone)]
struct AppState {
    vk_client: Arc<VkClient>,
    telegram_client: Arc<TelegramClient>,
    price_fetcher: Arc<PriceFetcher>,
}

async fn get_card_name(query: &str) -> Result<String> {
    let number_search_regex = Regex::new(r#"^!s (?<set>\w{3}) (?<number>\d+)$"#)
        .context("failed to compile number search regex")?;
    let name = if let Some(captures) = number_search_regex.captures(query) {
        let set = captures
            .name("set")
            .context("failed to get 'set' value from capture")?
            .as_str()
            .to_string();
        let number: usize = captures
            .name("number")
            .context("failed to get 'number' from capture")?
            .as_str()
            .parse()
            .context("failed to parse collector number")?;
        Card::set_and_number(&set, number)
            .await
            .map(|card| card.name)
            .context("failed to get card by set and number")?
    } else {
        Card::named_fuzzy(query)
            .await
            .map(|card| card.name)
            .context("failed to find card by it's fuzzy name")?
    };
    Ok(name)
}

#[derive(Serialize)]
struct StarCityRequestPayload {
    keyword: String,
    #[serde(rename = "ClientGuid")]
    client_guid: String,
    #[serde(rename = "SortBy")]
    sort_by: String,
    #[serde(rename = "FacetSelections")]
    facet_selections: StarCityFacetSelection,
}

#[derive(Serialize)]
struct StarCityFacetSelection {
    product_type: Vec<String>,
}

#[derive(Deserialize)]
struct StarCityResponse {
    #[serde(rename = "Results")]
    results: Vec<StarCityResponseResult>,
}

#[derive(Deserialize)]
struct StarCityResponseResult {
    #[serde(rename = "Document")]
    document: StarCityResponseDocument,
}

#[derive(Deserialize)]
struct StarCityResponseDocument {
    set: Vec<String>,
    hawk_child_attributes: Vec<StarCityResponseAttributesVariants>,
}

#[derive(Deserialize)]
#[serde(untagged)]
enum StarCityResponseAttributesVariants {
    Available(StarCityResponseAttributes),
    Unavailable {},
}

impl StarCityResponseAttributesVariants {
    fn get_card_info(&self, set: &str) -> Option<CardInfo> {
        if let StarCityResponseAttributesVariants::Available(res) = self {
            let condition = res.condition.first().cloned()?;
            let price = res.price.first().cloned()?.parse().ok()?;
            let url = res.url.first().cloned()?;
            if condition == "Near Mint" {
                Some(CardInfo {
                    set: set.to_string(),
                    price,
                    url: format!("https://starcitygames.com{url}"),
                })
            } else {
                None
            }
        } else {
            None
        }
    }
}

#[derive(Deserialize)]
struct StarCityResponseAttributes {
    price: Vec<String>,
    condition: Vec<String>,
    url: Vec<String>,
}

struct CardInfo {
    set: String,
    price: f32,
    url: String,
}

struct PriceFetcher {
    client_guid: String,
}

impl PriceFetcher {
    fn from_env() -> Result<Self> {
        let client_guid =
            env::var("SCG_CLIENT_GUID").context("SCG_CLIENT_GUID env variable is not set")?;
        Ok(Self { client_guid })
    }

    async fn get_card_prices(&self, name: &str) -> Result<Vec<CardInfo>> {
        let client = reqwest::ClientBuilder::new().build()?;
        let resp = client
            .post("https://essearchapi-na.hawksearch.com/api/v2/search")
            .json(&StarCityRequestPayload {
                keyword: name.to_string(),
                client_guid: self.client_guid.clone(),
                sort_by: "score".into(),
                facet_selections: StarCityFacetSelection {
                    product_type: vec!["Singles".into()],
                },
            })
            .send()
            .await
            .context("request to SCG failed")?;

        let response: StarCityResponse = resp.json().await.context("SCG returned invalid json")?;
        let res = response
            .results
            .iter()
            .flat_map(|result| {
                let set = result.document.set.first().cloned()?;
                let info = result
                    .document
                    .hawk_child_attributes
                    .iter()
                    .flat_map(|res| res.get_card_info(&set))
                    .collect::<Vec<_>>();
                Some(info)
            })
            .flatten()
            .collect::<Vec<_>>();
        Ok(res)
    }
}

struct TelegramClient {
    bot: Bot,
    secret: String,
}

impl TelegramClient {
    fn from_env() -> Result<Self> {
        let secret = env::var("TG_SECRET").context("failed to get TG_SECRET env var")?;
        Ok(Self {
            bot: Bot::from_env(),
            secret,
        })
    }
}

#[derive(Deserialize)]
struct TelegramUpdate {
    message: TelegramMessage,
}

#[derive(Deserialize)]
struct TelegramMessage {
    text: String,
    chat: TelegramChat,
}

#[derive(Deserialize)]
struct TelegramChat {
    id: i64,
}

async fn telegram(
    headers: HeaderMap,
    State(state): State<AppState>,
    Json(payload): Json<TelegramUpdate>,
) -> Result<&'static str, StatusCode> {
    let secret = headers
        .get("X-Telegram-Bot-Api-Secret-Token")
        .ok_or(StatusCode::FORBIDDEN)?
        .to_str()
        .map_err(|_| StatusCode::FORBIDDEN)?;
    if secret != state.telegram_client.secret {
        Err(StatusCode::FORBIDDEN)?;
    }
    let chat_id = payload.message.chat.id;
    let text = &payload.message.text;
    let res = handle_telegram_message(&state, chat_id, text).await;
    if let Err(err) = res {
        report_error(state.telegram_client.as_ref(), chat_id, err).await;
    }
    Ok("OK")
}

#[derive(Debug)]
enum ErrorContext {
    Scryfall,
    Scg,
    Telegram,
    Vk,
}

impl Display for ErrorContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let text = match self {
            Self::Scryfall => "failed to get card from scryfall",
            Self::Scg => "failed to get card from SCG",
            Self::Telegram => "failed to send message through Telergam",
            Self::Vk => "failed to send message through VK",
        };
        f.write_str(text)
    }
}

async fn handle_telegram_message(state: &AppState, chat_id: i64, message: &str) -> Result<()> {
    let name = get_card_name(message)
        .await
        .context(ErrorContext::Scryfall)?;
    let prices = state
        .price_fetcher
        .get_card_prices(&name)
        .await
        .context(ErrorContext::Scg)?;
    let content = prices
        .iter()
        .take(5)
        .enumerate()
        .map(|(i, info)| {
            format!(
                "{}\\. [{}]({}): {}",
                i + 1,
                escape(&info.set),
                escape_link_url(&info.url),
                escape(&format!("${}", info.price))
            )
        })
        .collect::<Vec<_>>()
        .join("\n");
    let header = escape(&format!("Оригинальное название: {}", name));
    let response = format!("{header}\n\n{content}");
    let request = SendMessage::new(ChatId(chat_id), response)
        .parse_mode(ParseMode::MarkdownV2)
        .disable_web_page_preview(true);
    JsonRequest::new(state.telegram_client.bot.clone(), request)
        .send()
        .await
        .context(ErrorContext::Telegram)?;
    Ok(())
}

struct VkClient {
    client: reqwest::Client,
    token: String,
    group_id: i64,
    confirmation: String,
    secret: String,
}

impl VkClient {
    fn from_env() -> Result<Self> {
        let client = reqwest::Client::new();
        let token = env::var("VK_TOKEN").context("failed to get VK_TOKEN")?;
        let group_id = env::var("VK_GROUP_ID")
            .context("failed to get VK_GROUP_ID")
            .and_then(|x| x.parse().context("failed to parse VK_GROUP_ID as a number"))?;
        let confirmation =
            env::var("VK_CONFIRMATION_STRING").context("failed to get VK_CONFIRMATION_STRING")?;
        let secret = env::var("VK_SECRET").context("failed to get VK_SECRET")?;
        Ok(Self {
            client,
            token,
            group_id,
            confirmation,
            secret,
        })
    }

    async fn send(&self, user_id: i64, message: &str) -> Result<()> {
        self.client
            .get("https://api.vk.com/method/messages.send")
            .query(&[
                ("user_id", user_id.to_string().as_str()),
                ("v", "5.131"),
                ("access_token", &self.token),
                ("random_id", "0"),
                ("message", message),
            ])
            .send()
            .await
            .context(ErrorContext::Vk)?;
        Ok(())
    }
}

#[derive(Deserialize)]
#[serde(tag = "type")]
enum VkRequest {
    #[serde(rename = "message_new")]
    Message(VkMessageRequest),
    #[serde(rename = "confirmation")]
    Confirmation(VkConfirmationRequest),
}

#[derive(Deserialize)]
struct VkMessageRequest {
    object: VkMessageObject,
    secret: String,
}

#[derive(Deserialize)]
struct VkMessageObject {
    from_id: i64,
    text: String,
}

#[derive(Deserialize)]
struct VkConfirmationRequest {
    group_id: i64,
}

async fn vk(State(state): State<AppState>, Json(request): Json<VkRequest>) -> (StatusCode, String) {
    match request {
        VkRequest::Message(payload) => {
            if payload.secret != state.vk_client.secret {
                return (StatusCode::FORBIDDEN, "Access denied".into());
            }
            let user_id = payload.object.from_id;
            let message = &payload.object.text;
            let res = handle_vk_message(&state, user_id, message).await;
            if let Err(err) = res {
                report_error(state.vk_client.as_ref(), user_id, err).await;
            }
        }
        VkRequest::Confirmation(confirmation) => {
            if confirmation.group_id != state.vk_client.group_id {
                return (StatusCode::FORBIDDEN, "Access denied".into());
            }
            return (StatusCode::OK, state.vk_client.confirmation.clone());
        }
    }
    (StatusCode::OK, "OK".into())
}

async fn handle_vk_message(state: &AppState, user_id: i64, message: &str) -> Result<()> {
    let name = get_card_name(message)
        .await
        .context(ErrorContext::Scryfall)?;
    let prices = state
        .price_fetcher
        .get_card_prices(&name)
        .await
        .context(ErrorContext::Scg)?;
    let header = escape(&format!("Оригинальное название: {}", name));
    let content = prices
        .iter()
        .take(5)
        .enumerate()
        .map(|(i, info)| format!("{}. {}: ${}\n{}", i + 1, info.set, info.price, info.url))
        .collect::<Vec<_>>()
        .join("\n");
    let response = format!("{header}\n\n{content}");
    state.vk_client.send(user_id, &response).await
}

trait MessageSender {
    async fn send(&self, user_id: i64, message: &str) -> Result<()>;
}

impl MessageSender for VkClient {
    async fn send(&self, user_id: i64, message: &str) -> Result<()> {
        self.send(user_id, message).await?;
        Ok(())
    }
}

impl MessageSender for TelegramClient {
    async fn send(&self, user_id: i64, message: &str) -> Result<()> {
        self.bot.send_message(ChatId(user_id), message).await?;
        Ok(())
    }
}

async fn report_error<T: MessageSender>(sender: &T, chat_id: i64, err: anyhow::Error) {
    if !matches!(
        err.downcast_ref::<scryfall::Error>(),
        Some(scryfall::Error::ScryfallError(_))
    ) {
        println!("error: {:#}", err);
    }
    let sent = match err.downcast_ref::<ErrorContext>() {
        Some(ErrorContext::Scryfall) => sender
            .send(chat_id, "Карта не найдена")
            .await
            .map(|_| ())
            .context("failed to send error message"),
        Some(ErrorContext::Scg) => sender
            .send(chat_id, "Цены не найдены")
            .await
            .map(|_| ())
            .context("failed to send error message"),
        _ => Ok(()),
    };
    if let Err(err) = sent {
        println!("error: {:#}", err);
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let port = env::var("PORT")
        .ok()
        .and_then(|x| x.parse().ok())
        .unwrap_or(3000);
    let addr: (&str, u16) = ("0.0.0.0", port);
    let listener = TcpListener::bind(addr)
        .await
        .context("failed to create tcp listener")?;
    let vk_client = VkClient::from_env().context("failed to init vk client")?;
    let telegram_client = TelegramClient::from_env().context("failed to init telegram client")?;
    let price_fetcher = PriceFetcher::from_env().context("failed to init price fetcher")?;
    let state = AppState {
        vk_client: Arc::new(vk_client),
        telegram_client: Arc::new(telegram_client),
        price_fetcher: Arc::new(price_fetcher),
    };
    let app = Router::new()
        .route("/tg", post(telegram))
        .route("/vk", post(vk))
        .with_state(state);
    axum::serve(listener, app).await.unwrap();
    Ok(())
}
