import {ChangeDetectionStrategy, Component, Input, OnInit} from '@angular/core';
import {marked} from 'marked';
import {DomSanitizer} from "@angular/platform-browser";

@Component({
    selector: 'app-markdown',
    templateUrl: './markdown.component.html',
    styleUrls: ['./markdown.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class MarkdownComponent implements OnInit {

  @Input()
  set raw(value: string) {
    this._raw = value;
    this.render();
  }

  get raw(): string {
    return this._raw;
  }

  private _raw: string;

  renderedHTML: any;

  constructor(public sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
  }

  render() {

    const renderer = new marked.Renderer();

    /**
     * Make Marked support specifying image size in pixels in this format:
     *
     * ![alt](src = x WIDTH)
     * ![alt](src = HEIGHT x)
     * ![alt](src = HEIGHT x WIDTH)
     * ![alt](src = x WIDTH "title")
     * ![alt](src = HEIGHT x "title")
     * ![alt](src = HEIGHT x WIDTH "title")
     *
     * Note: whitespace from the equals sign to the title/end of image is all
     * optional. Each of the above examples are equivalent to these below,
     * respectively:
     *
     * ![alt](src =xWIDTH)
     * ![alt](src =HEIGHTx)
     * ![alt](src =HEIGHTxWIDTH)
     * ![alt](src =xWIDTH "title")
     * ![alt](src =HEIGHTx "title")
     * ![alt](src =HEIGHTxWIDTH "title")
     *
     * Example usage:
     *
     * ![my image](https://example.com/my-image.png =400x600 "My image")
     * ![](https://example.com/my-image.png =400x "My image")
     * ![](https://example.com/my-image.png =400x)
     */
    renderer.image = function (src, title, alt) {
      const parts = /(.*)\s+=\s*(\d*)\s*x\s*(\d*)\s*$/.exec(src)
      var url = src
      var height = undefined
      var width = undefined
      if (parts) {
        url = parts[1]
        height = parts[2]
        width = parts[3]
      }
      var YouTube = mediaParseIdFromUrl('youtube', url);
      var Vimeo = mediaParseIdFromUrl('vimeo', url);
      var Viddler = mediaParseIdFromUrl('viddler', url);
      var DailyMotion = mediaParseIdFromUrl('dailymotion', url);
      var Html5 = mediaParseIdFromUrl('html5', url);
      let res = ''
      if (YouTube !== undefined) {
        res = create_iframe('//www.youtube.com/embed/' + YouTube, title, alt, height, width);
      } else if (Vimeo !== undefined) {
        res = create_iframe('//player.vimeo.com/video/' + Vimeo + '?api=1', title, alt, height, width);
      } else if (Viddler !== undefined) {
        res = create_iframe('//www.viddler.com/player/' + Viddler, title, alt, height, width);
      } else if (DailyMotion !== undefined) {
        res = create_iframe('//www.dailymotion.com/embed/video/' + DailyMotion, title, alt, height, width);
      } else if (Html5) {
        res = '<video';
        if (height) res += ' height="' + height + '"'
        if (width) res += ' width="' + width + '"'
        res += ' controls><source src="' + Html5['link'] + '" type="video/' + Html5['extension'] + '">'
        if (alt) res += sanitize(alt)
        res += '</video>';
      } else {
        res = '<img '
        if (height) res += ' height="' + height + '"'
        if (width) res += ' width="' + width + '"'
        res += 'src="' + sanitize(url) + '"'
        if (alt) res += ' alt="' + sanitize(alt) + '"'
        if (title) res += ' title="' + sanitize(title) + '"'
        res += '>'
      }
      return res
    }
    let markd = marked.setOptions({
      renderer: renderer
    })
    this.renderedHTML = markd(this.raw);
  }
}


function sanitize(str) {
  return str.replace(/&<"/g, function (m) {
    if (m === "&") return "&amp;"
    if (m === "<") return "&lt;"
    return "&quot;"
  })
}

/**
 * Parse url of video to return Video ID only
 * if video exists and matches to media's host
 * else undefined
 *
 * @example mediaParseIdFromUrl('youtube', 'https://www.youtube.com/watch?v=fgQRaRqOTr0')
 * //=> fgQRaRqOTr0
 *
 * @param  {string} provider    name of media/video site
 * @param  {string} url         url of video
 * @return {string|undefined}   the parsed id of video, if not match - undefined
 */
function mediaParseIdFromUrl(provider, url) {
  if (provider === 'youtube') {
    var youtubeRegex = /^.*((youtu.be\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?]*).*/;
    var youtubeMatch = url.match(youtubeRegex);
    if (youtubeMatch && youtubeMatch[7].length == 11) {
      return youtubeMatch[7];
    } else {
      return undefined;
    }
  } else if (provider === 'vimeo') {
    var vimeoRegex = /^.*vimeo.com\/(\d+)/;
    var vimeoMatch = url.match(vimeoRegex);
    if (vimeoMatch && vimeoMatch[1].length == 8) {
      return vimeoMatch[1];
    } else {
      return undefined;
    }
  } else if (provider === 'viddler') {
    var viddlerRegex = /^.*((viddler.com\/)|(v\/)|(\/u\/\w\/)|(embed\/)|(watch\?))\??v?=?([^#\&\?]*).*/;
    var viddlerMatch = url.match(viddlerRegex);
    if (viddlerMatch && viddlerMatch[7].length == 8) {
      return viddlerMatch[7];
    } else {
      return undefined;
    }
  } else if (provider === 'dailymotion') {
    var dailymotionRegex = /^.+dailymotion.com\/((video|hub)\/([^_]+))?[^#]*(#video=([^_&]+))?/;
    var dailymotionMatch = url.match(dailymotionRegex);
    if (dailymotionMatch && (dailymotionMatch[5] || dailymotionMatch[3])) {
      if (dailymotionMatch[5]) {
        return dailymotionMatch[5];
      }
      if (dailymotionMatch[3]) {
        return dailymotionMatch[3];
      }
      return undefined;
    } else {
      return undefined;
    }
  } else if (provider === 'html5') {
    var html5Regex = /.+\.(wav|mp3|ogg|mp4|wma|webm|mp3)$/i;
    var html5Match = url.match(html5Regex);

    if (html5Match && html5Match[1]) {
      return {'link': html5Match[0], 'extension': html5Match[1]};
    } else {
      return undefined;
    }
  } else {
    return undefined;
  }
}

function create_iframe(src, title, alt, height, width) {
  var res = '<iframe';
  if (title) res += ' title="' + sanitize(title) + '"'
  if (height) res += ' height="' + height + '"'
  if (width) res += ' width="' + width + '"'
  res += ' src="' + src + '" frameborder="0" webkitAllowFullScreen mozallowfullscreen allowFullScreen>';
  res += sanitize(alt) + '</iframe>';
  return res;
}

