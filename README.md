# MeiliSitemap
MeiliSitemap is sitemap generator for Meilisearch indexes and support normal, video, news and image structure.

## Features

- Create sitemap from multiple index
- Support gzip compression
- Local file server for sitemap
- Support custom name for sitemaps (default is index name)
- Support sitemap stylesheets
- Support custom path for sitemaps
- Support filters for get specific documents
- Support live update sitemap with background scheduler
- Support normal, video, image and news sitemap type

## Installation

You can install `meilisitemap` using different methods.

### Release

Download the latest release from [here](https://github.com/Ja7ad/meilisitemap/releases).

### Go Installation

Install the package using Go:

```shell
go install github.com/Ja7ad/meilisitemap/cmd/meilisitemap@latest
```

### Docker

Run `meilisitemap` with real-time sync using Docker:

- Docker

```shell
docker run --rm -it -v ./config.yml:/etc/meilisitemap/config.yml --name meilisitemap ja7adr/meilisitemap
```

- Docker compose

```yaml
version: "3"
services:
  meilibridge:
    image: ja7adr/meilisitemap:latest
    volumes:
      - ./config.yml:/etc/meilisitemap/config.yml
    restart: always
```

## Example Configuration

example configuration for run meilisitemap

```yaml
general:
  # base url is your site base url for document link or loc (require)
  base_url: https://example.com
  # indexsitemap_base_url set sitemap link of each index and put in sitemap.xml file as sitemapindex.
  # base_url + indexsitemap_path = index_sitemap_link
  # for example https://example.com/sitemaps/movies.xml
  # note: if serve enabled set indexsitemap_base_url to serve.listen/sitemaps/
  indexsitemap_path: /sitemaps/
  # set custom name for sitemap file
  # default is null and sitemap.xml
  file_name: meilisearch_sitemap.xml
  # set prefix for all sitemaps
  # for example in result meilisearch_sitemap.xml
  # default is null
  prefix: "meilisearch_"
  # set custom stylesheet for sitemap, currently style1 and style2 is available
  # default is null
  stylesheet: style1
  # available your sitemap on local server
  # for example http://127.0.0.1:8080/sitemap.xml
  # note1: if serve enable, possible your sitemapindex urls set to http://127.0.0.1:8080/sitemaps/movies.xml
  # note2: if you don't set indexsitemap_base_url, auto enable serve.
  serve:
    enable: true
    listen: 127.0.0.1:8080
    pprof: false

  # meilisearch host and api_key (require)
  meilisearch:
    host: "http://localhost:7700"
    api_key: "masterKey"


# sitemaps create specific sitemap file for every index and put in sitemap.xml as sitemapindex (require)
sitemaps:
  # index name
  movies:
    # make xml sitemap (require)
    sitemap: true
    # make html sitemap
    html_sitemap: true
    # make rss feed
    rss: false
    # meilisearch filter expression
    # https://www.meilisearch.com/docs/learn/filtering_and_sorting/filter_expression_reference
    # default is null and make sitemap for all documents
    filter: "genre = horror & imdb_rate > 5"
    # base path is document group path, final address base_url + base_path + unique_field = loc or link (require)
    base_path: "/movies/"
    # compress with gzip
    # for example result is movies.gz
    compress: false
    # set custom name for index sitemap file
    # default is null and index name for file.
    sitemap_file_name: foobar
    # auto update sitemap in background by scheduler, duration is base on changefreq
    live_update:
      enabled: false
      # interval scheduler duration in seconds
      interval: 3000

    # map document fields to sitemap structure (require)
    field_map:
      # unique_field final things for loc you can set title, name, id or unique number
      # don't set long text similar description.
      # for text field auto replace - with space, example for text (title: Anatomy of a Fall -> anatomy-of-a-fall)
      # (require)
      unique_field: title
      # lastmod is W3C date and time format.
      # if you don't have date-time field auto set current datetime.
      lastmod: created_at
      # changefreq: always, hourly, daily, weekly, monthly, yearly, never
      # default is daily
      changefreq: "daily"
      # priority: low (0.3), medium (0.5), high (0.8), highest (1.0)
      # default is high
      priority: high
      # Optional video field map for movies that include video data
      video:
          # URL to the video thumbnail
          # if you have file id or file name for thumbnail you can set base url for file id with extension if required
          # for example: "image_id|https://cdn.example.com/images|.jpg"
          thumbnail_loc: thumbnail_url
          title: title                  # Title of the video
          description: description      # Description of the video
          # URL to the video content
          # if you have file id or file name for video you can set base url for video id
          # for example: "video_id|https://cdn.example.com/videos"
          content_loc: video_url
          player_loc: player_url        # URL to the video player
          duration: duration            # Duration of the video in seconds
          expiration_date: expiration   # Expiration date of the video
          rating: rating                # Rating of the video
          view_count: view_count        # View count of the video
          publication_date: published_at # Publication date of the video
          family_friendly: family_friendly # Family-friendly flag
          restriction: restriction      # Restriction details for the video
          requires_subscription: requires_subscription # Subscription requirement flag
          live: live                    # Live broadcast flag
      # Optional image field map for movies that include image data
      image:
        # URL to the image
        # if you have file id or file name for image you can set base url for file id
        # for example: "image_id|https://cdn.example.com/images|.png"
        loc: image_url
        caption: image_caption        # Caption for the image
        title: image_title            # Title of the image
        license: image_license        # License for the image
        geo_location: image_location  # Geolocation for the image
      # Optional news field map for movies that include news data
      news:
        publication:
          name: publication_name      # Name of the publication
          language: publication_language # Language of the publication
        pub_date: publication_date    # Publication date of the news
        title: news_title             # Title of the news article
        keywords: news_keywords       # Keywords for the news article
        description: news_description # Description of the news article

  category:
    # make xml sitemap
    sitemap: true
    # make html sitemap
    html_sitemap: true
    # make rss feed
    rss: false
    # base path is item address base_url + base_path + unique_field = loc or link
    base_path: "/categories/"
    # compress with gzip
    compress: true

    field_map:
      # unique_field final things for loc you can set title, name, id or unique number
      # don't set long text similar description.
      # for text field auto replace - with space, example for text (title: Anatomy of a Fall -> anatomy-of-a-fall)
      unique_field: id
      # lastmod is W3C date and time format.
      lastmod: created_at
      # changefreq: always, hourly, daily, weekly, monthly, yearly, never
      changefreq: "daily"
      # priority: low (0.3), medium (0.5), high (0.8), highest (1.0)
      # default is high
      priority: medium
      # Optional video field map for movies that include video data
      video:
        thumbnail_loc: thumbnail_url  # URL to the video thumbnail
        title: title                  # Title of the video
        description: description      # Description of the video
        content_loc: video_url        # URL to the video content
        player_loc: player_url        # URL to the video player
        player_auto_play: autoplay    # player_auto_play is attribute for player_loc ap=1 or 0
        duration: duration            # Duration of the video in seconds
        expiration_date: expiration   # Expiration date of the video
        rating: rating                # Rating of the video
        view_count: view_count        # View count of the video
        publication_date: published_at # Publication date of the video
        family_friendly: family_friendly # Family-friendly flag
        restriction: restriction      # Restriction details for the video
        relationship: restriction_relationship # Relationship attribute for restriction allow or deny
        requires_subscription: requires_subscription # Subscription requirement flag
        live: live                    # Live broadcast flag
      # Optional image field map for movies that include image data
      image:
        loc: image_url                # URL to the image
        caption: image_caption        # Caption for the image
        title: image_title            # Title of the image
        license: image_license        # License for the image
        geo_location: image_location  # Geolocation for the image
      # Optional news field map for movies that include news data
      news:
        publication:
          name: publication_name      # Name of the publication
          language: publication_language # Language of the publication
        pub_date: publication_date    # Publication date of the news
        title: news_title             # Title of the news article
        keywords: news_keywords       # Keywords for the news article
        description: news_description # Description of the news article
```