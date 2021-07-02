import marked, { Renderer } from 'marked';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';

import {ChangeDetectionStrategy, Component, Input, OnChanges, SimpleChange, ViewEncapsulation} from '@angular/core';
import {DomSanitizer, SafeHtml} from "@angular/platform-browser";

@Component({
  selector: 'ktb-markdown',
  templateUrl: './ktb-markdown.component.html',
  styleUrls: ['./ktb-markdown.component.scss'],
  host: {
    class: 'ktb-markdown'
  },
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbMarkdownComponent implements OnChanges {

  @Input() markdown: string;
  @Input() html: string;
  public safeHtml: SafeHtml;
  private md: any;

  static highlightCode(code: string, language: string): string {
    if (!(language && hljs.getLanguage(language))) {
      // use 'markdown' as default language
      language = 'markdown';
    }

    const result = hljs.highlight(language, code).value;
    return `<code class="hljs ${language}">${result}</code>`;
  }

  static addTargetAndNoopener(node) {
    // set all elements owning href to target=_blank and rel=noopener
    if ('href' in node) {
      node.setAttribute('target', '_blank');
      node.setAttribute('rel', 'noopener');
    }
  }

  constructor(private sanitizer: DomSanitizer) {
    const renderer = new Renderer();
    renderer.code = KtbMarkdownComponent.highlightCode;
    DOMPurify.addHook('afterSanitizeAttributes', KtbMarkdownComponent.addTargetAndNoopener);
    this.md = marked.setOptions({ renderer });
  }

  markdownToSafeHtml(value: string): SafeHtml {
    const html = this.md(value);
    const safeHtml = DOMPurify.sanitize(html);
    return this.sanitizer.bypassSecurityTrustHtml(safeHtml);
  }

  htmlToSafeHtml(value: string): SafeHtml {
    const safeHtml = DOMPurify.sanitize(value);
    return this.sanitizer.bypassSecurityTrustHtml(safeHtml);
  }

  ngOnChanges(changes: { [propKey: string]: SimpleChange }) {
    for (const propName in changes) {
      if (propName === 'markdown') {
        const value = changes[propName].currentValue;
        if (value) {
          this.safeHtml = this.markdownToSafeHtml(value);
        }
      } else if (propName === 'html') {
        const value = changes[propName].currentValue;
        if (value) {
          this.safeHtml = this.htmlToSafeHtml(value);
        }
      }
    }
  }
}
