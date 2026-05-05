const DESKTOP_OPEN_URL_MESSAGE_PREFIX = "mistermorph:open-url:";

function currentWindow() {
  return typeof window === "undefined" ? null : window;
}

function resolveHTTPURL(rawURL) {
  const value = String(rawURL || "").trim();
  if (!value) {
    return "";
  }
  const win = currentWindow();
  const baseURL = win?.location?.href || "http://localhost/";
  try {
    const url = new URL(value, baseURL);
    if (url.protocol !== "http:" && url.protocol !== "https:") {
      return "";
    }
    return url.href;
  } catch {
    return "";
  }
}

function desktopMessageSender() {
  const win = currentWindow();
  if (!win) {
    return null;
  }

  const chromePostMessage = win.chrome?.webview?.postMessage;
  if (typeof chromePostMessage === "function") {
    return (message) => chromePostMessage.call(win.chrome.webview, message);
  }

  const webkitPostMessage = win.webkit?.messageHandlers?.external?.postMessage;
  if (typeof webkitPostMessage === "function") {
    return (message) => webkitPostMessage.call(win.webkit.messageHandlers.external, message);
  }

  const wailsInvoke = win._wails?.invoke || win.wails?.invoke;
  if (typeof wailsInvoke === "function") {
    return (message) => wailsInvoke(message);
  }

  return null;
}

function postDesktopMessage(message) {
  const send = desktopMessageSender();
  if (!send) {
    return false;
  }
  try {
    send(message);
    return true;
  } catch {
    return false;
  }
}

export function canOpenExternalURLInDesktop() {
  return desktopMessageSender() !== null;
}

export function openExternalPlaceholder() {
  const win = currentWindow();
  if (!win || typeof win.open !== "function") {
    return null;
  }
  try {
    const popup = win.open("about:blank", "_blank");
    if (popup) {
      popup.opener = null;
    }
    return popup;
  } catch {
    return null;
  }
}

export function openExternalURL(rawURL) {
  const target = resolveHTTPURL(rawURL);
  if (!target) {
    return false;
  }

  if (postDesktopMessage(`${DESKTOP_OPEN_URL_MESSAGE_PREFIX}${target}`)) {
    return true;
  }

  const win = currentWindow();
  if (!win || typeof win.open !== "function") {
    return false;
  }
  const popup = win.open(target, "_blank", "noopener,noreferrer");
  if (popup) {
    popup.opener = null;
  }
  return true;
}

function isSameOriginURL(target) {
  const win = currentWindow();
  if (!win?.location?.origin) {
    return false;
  }
  try {
    return new URL(target).origin === win.location.origin;
  } catch {
    return false;
  }
}

function handleDocumentClick(event) {
  if (!canOpenExternalURLInDesktop() || event.defaultPrevented || event.button !== 0) {
    return;
  }
  const eventTarget = event.target instanceof Element ? event.target : event.target?.parentElement;
  const anchor = eventTarget?.closest?.("a[href]");
  if (!anchor) {
    return;
  }
  const target = resolveHTTPURL(anchor.getAttribute("href"));
  if (!target || isSameOriginURL(target)) {
    return;
  }
  event.preventDefault();
  openExternalURL(target);
}

export function installExternalLinkHandler() {
  if (typeof document === "undefined") {
    return;
  }
  document.addEventListener("click", handleDocumentClick);
}
