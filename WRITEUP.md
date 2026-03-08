## Motivation and Background

During the early stage of crawler development, our team frequently relied on third-party scraping infrastructure such as services from companies like Bright Data and ZenRows because many target websites were heavily protected by advanced WAF systems like Akamai, Cloudflare, and DataDome.

Our internal crawlers were initially built using traditional browser automation frameworks such as Playwright and Selenium. These tools are extremely convenient for automation, but they expose many signals that allow anti-bot systems to detect automation.

Examples of detectable signals included:

* `navigator.webdriver = true`
* injected Playwright objects like `window.__playwright__binding__`
* injected JavaScript execution environments
* mismatches between browser properties and network headers
* detectable overrides through `Object.getOwnPropertyDescriptor`
* hijacked functions where `toString()` no longer returns `[native code]`

Another major problem is that browser properties are spread across multiple execution contexts: main world, worker threads, network stack, GPU stack, etc. JavaScript-based spoofing only modifies the page context, leaving deeper layers untouched.

Because of this, spoofing through JavaScript injections creates **inconsistent fingerprints**. Modern anti-bot systems are designed to detect exactly these inconsistencies.

For example:

* network User-Agent ≠ `navigator.userAgent`
* OS mismatch with GPU renderer
* timezone inconsistent with IP region
* browser APIs behaving differently across contexts

Even stealth scripts injected through Playwright initialization hooks could only partially hide automation, and were often detected by deeper inspection.

This led me to investigate **why a real user browser passes detection while automated browsers fail**, even when they appear similar on the surface.

After extensive research, the key insight was:

> The problem is not automation itself — the problem is fingerprint inconsistency across layers.

Automation frameworks simply do not allow modification of many internal browser signals because they operate outside the browser engine.

---

## Building a Patched Browser (Camoium)

To solve this problem, I started modifying the browser itself.

I chose to work from Ungoogled Chromium because it already removes dependencies on Google services and is easier to modify.

From there I built a patched Chromium variant called **Camoium**.

The idea was simple:

Instead of spoofing properties through JavaScript, modify them **directly inside the browser’s native implementation (C++ level)**.

Camoium intercepts browser API calls internally so spoofed values appear as **native browser values**, not overwritten JavaScript properties.

This eliminates many detection techniques used by WAFs.

Key characteristics:

* fingerprint injected at browser startup
* no runtime JavaScript spoofing
* consistent identity across browser subsystems
* native browser APIs return spoofed values

---

## Automation Layer (undetected-cdp)

To control Camoium, I built a CDP-based automation tool similar to Playwright called **undetected-cdp**.

It generates a realistic device identity and passes it to the browser at startup.

That identity includes:

* OS
* hardware profile
* screen resolution
* GPU profile
* fonts
* timezone
* locale
* WebRTC configuration
* browser version
* network headers

Because the browser itself applies the fingerprint internally, every API returns **consistent values across all contexts**.

This approach significantly improves success rates compared to JavaScript-based stealth patches.

---

## Fingerprint Coverage

Camoium modifies or spoofs many high-entropy fingerprint signals including:

**Navigator properties**

* device
* OS
* browser version
* hardware concurrency
* memory

**Display environment**

* screen size
* viewport
* window dimensions
* devicePixelRatio

**Localization**

* timezone
* geolocation
* language
* Intl APIs

**Graphics stack**

* WebGL parameters
* GPU renderer
* supported extensions
* shader precision formats

**Media stack**

* AudioContext sample rate
* output latency
* max channels
* speech voices

**Network layer**

* User-Agent
* Accept-Language
* header order

**Other APIs**

* Battery API
* WebRTC IP
* font enumeration

Additionally, Camoium patches headless mode so that it behaves identically to a normal windowed browser.

---

## Success Rate Factors

In practice, success rate depends mainly on three factors:

### 1. Fingerprint Consistency

All parts of the fingerprint must logically match.

Examples of impossible combinations:

* Windows user agent + Apple GPU
* macOS user agent + DirectX renderer
* mobile user agent + desktop resolution

Modern anti-bot systems immediately flag these inconsistencies.

### 2. Proxy Quality

IP reputation is still important.

Rotating IPs alone does not work if the browser fingerprint stays identical.

Each browser instance needs a **unique identity**, otherwise rate-limiting or fingerprint correlation will reveal automation.

### 3. Request Behavior

Automation patterns such as navigation timing, interaction speed, and request burst patterns can also trigger detection.

---

## Desktop vs Mobile Differences

Mobile emulation is significantly more complex than desktop spoofing.

Desktop browsers expose many signals that do not exist on mobile devices.

Key differences include:

* touch input capabilities (`maxTouchPoints`)
* devicePixelRatio
* viewport scaling
* mobile-specific APIs
* sensor availability
* GPU and WebGL differences
* network header ordering
* canvas and audio fingerprint characteristics

Although Camoium can generate mobile fingerprints, mobile spoofing from desktop hardware is still more difficult to perfect.

Some WAFs still detect automation likely due to differences in **input event handling or interaction models**.

Further research is needed to fully replicate native mobile behavior.

---

## Regional Differences

Some anti-bot systems apply different detection thresholds depending on geographic region.

Examples include:

* stricter checks for high-fraud regions
* IP reputation scoring differences
* localized content affecting DOM behavior
* different CAPTCHA triggers

Additionally, region must align with fingerprint properties such as:

* timezone
* language
* geolocation
* IP country

If these do not match, detection probability increases.

---

# 2. Question Answers (Anti-Bot Engineer: Take-Home Challenge)

## Your approach and why you chose it

Traditional automation tools rely on JavaScript injection to hide automation signals, but this approach creates inconsistencies across browser layers. Anti-bot systems are designed to detect these inconsistencies.

To solve this, I modified the browser itself rather than patching behavior externally.

I built a patched Chromium browser called `Camoium` that injects the fingerprint directly at the native implementation level. This allows all browser APIs to return consistent values without JavaScript overrides, making the fingerprint appear fully native.

---

## What detection signals you encountered and how you identified them

The most common detection signals included:

* `navigator.webdriver`
* Playwright injected objects in the window scope
* mismatches between network headers and browser APIs
* detectable property overrides via `Object.getOwnPropertyDescriptor`
* hijacked functions where `toString()` reveals modifications
* inconsistencies between main thread and worker thread environments
* GPU and WebGL fingerprint mismatches

These signals were identified by comparing behavior between automated browsers and real user browsers and by analyzing how WAFs responded to different fingerprint configurations.

---

## What you tried that didn't work and why

Initially I tried stealth JavaScript injection techniques to patch browser properties.

However this approach had several limitations:

* JavaScript cannot modify deeper browser layers such as GPU or network stack
* overwritten properties are detectable
* injected scripts create inconsistent fingerprints
* worker contexts expose unpatched values

Because of these limitations, JavaScript-based stealth solutions were unreliable against advanced anti-bot systems.

---

## What affects the success rate — what matters and what doesn't

The most important factors are:

1. fingerprint consistency across all browser layers
2. proxy quality and IP reputation
3. realistic browser behavior patterns

What matters less is simple IP rotation.

If the fingerprint remains identical across requests, anti-bot systems can correlate activity regardless of IP changes.

A believable device identity combined with realistic browsing behavior has a much larger impact on success rate.

---

## Differences between desktop and mobile

Mobile environments expose different fingerprint signals such as touch capabilities, viewport scaling, and device-specific graphics behavior.

Emulating mobile devices from desktop hardware is challenging because some signals originate from physical hardware characteristics.

For this reason mobile spoofing often requires additional work beyond standard desktop fingerprint rotation.

---

## Differences across regions

Detection behavior can vary by region due to IP reputation scoring and fraud risk models.

Some regions trigger stricter verification or CAPTCHA challenges.

Fingerprint properties such as timezone, language, and geolocation must align with the proxy region to maintain credibility.

---

## What you would improve with more time

With more time I would focus on improving mobile emulation, particularly around input events and high-entropy signals that differ from desktop environments.

I would also expand fingerprint diversity by modeling more real-world device profiles and collecting additional telemetry from real browsers to better replicate natural behavior patterns.
