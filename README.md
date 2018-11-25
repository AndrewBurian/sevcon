# Sevcon
Turns your Pager Duty feed into your own DEFCON style warning system.

## Running
If not running in test mode, the only required config is a `PAGERDUTY_TOKEN` environment variable. To get an API key see the [PagerDuty Docs](https://support.pagerduty.com/docs/using-the-api#section-generating-an-api-key). The key can and should be read-only.

Other flags available:
- `-debug`: log at a higher verbosity
- `-quiet`: log errors only
- `-port <n>`: bind to the specific port (default 8080)
- `-test`: don't call the PD API, and instead just rotate through levels
