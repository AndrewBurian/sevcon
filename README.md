# SEVCON
Turns your Pager Duty feed into your own DEFCON style warning system.

SEVCON polls your PagerDuty accounts looking for incidents with assigned priorities. If the priority tag is in the form `SEV-[1-5]` then SEVCON will report the highest priority SEV (1 high, 5 low) and set the SEVCON light appropriately.

The intended use is to run on a communal display and bring a little more wide visibility into the current state of your incident response team's lives.

## Running
If not running in test mode, the only required config is a `PAGERDUTY_TOKEN` environment variable. To get an API key see the [PagerDuty Docs](https://support.pagerduty.com/docs/using-the-api#section-generating-an-api-key). The key can and should be read-only.

Other flags available:
- `-debug`: log at a higher verbosity
- `-quiet`: log errors only
- `-port <n>`: bind to the specific port (default 8080)
- `-test`: don't call the PD API, and instead just rotate through levels
- `-rate`: adjust the polling rate (default 1m)

## Web Interface
SEVCON runs a web interface at the root of the project `/` or `/index.html`.

## API
The only API is an HTTP SSE connection available at `/updates`. The web interface uses this as well.

No authentication is required. Handle with care.

Connecting to the updates feed should result in the latest update being delivered immediately. Further updates are delivered as the system state changes. No updates are delivered if a new SEV is opened with lower or the same priority as the current state.