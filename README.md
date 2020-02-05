# wantarr

A simple CLI tool that can be used to search for wanted media in:

- Sonarr
- Radarr

Once an item has been searched, it will not be searched again until the retry days setting has been reached.

## Configuration

```yaml
pvr:
  sonarr:
    type: sonarr
    url: https://sonarr.domain.com
    api_key: YOUR_API_KEY
    retry_days_age:
      missing: 90
      cutoff: 90
  radarr:
    type: radarr
    url: https://radarr.domain.com
    api_key: YOUR_API_KEY
    retry_days_age:
      missing: 90
      cutoff: 90
  radarr4k:
    type: radarr
    url: https://radarr.domain.com
    api_key: YOUR_API_KEY
    retry_days_age:
      missing: 90
      cutoff: 90
```


## Examples

- `wantarr missing radarr -v -m 20`
- `wantarr cutoff radarr4k -v -m 20`

## Notes

Supported Sonarr Version(s):

- 3

Supported Radarr Version(s):

- 0.2