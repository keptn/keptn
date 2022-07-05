import { marked, Renderer } from 'marked';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';
import {
  ChangeDetectionStrategy,
  Component,
  HostBinding,
  Input,
  OnChanges,
  SimpleChange,
  ViewEncapsulation,
} from '@angular/core';
import { SafeHtml } from '@angular/platform-browser';

@Component({
  selector: 'ktb-markdown',
  templateUrl: './ktb-markdown.component.html',
  styleUrls: ['./ktb-markdown.component.scss'],
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbMarkdownComponent implements OnChanges {
  @HostBinding('class') cls = 'ktb-markdown';
  @Input() markdown?: string;
  @Input() html?: string;
  public safeHtml?: SafeHtml;
  private readonly md: typeof marked;

  static highlightCode(code: string, language: string): string {
    if (!(language && hljs.getLanguage(language))) {
      // use 'markdown' as default language
      language = 'markdown';
    }

    const result = hljs.highlight(language, code).value;
    return `<code class="hljs ${language}">${result}</code>`;
  }

  static addTargetAndNoopener(node: Element): void {
    // set all elements owning href to target=_blank and rel=noopener
    if (node instanceof HTMLLinkElement) {
      node.setAttribute('target', '_blank');
      node.setAttribute('rel', 'noopener');
    }
  }

  constructor() {
    const renderer = new Renderer();
    renderer.code = KtbMarkdownComponent.highlightCode;
    DOMPurify.addHook('afterSanitizeAttributes', KtbMarkdownComponent.addTargetAndNoopener);
    this.md = marked.setOptions({ renderer });
  }

  markdownToSafeHtml(value: string): SafeHtml {
    const html = this.md(value);
    return DOMPurify.sanitize(html);
  }

  htmlToSafeHtml(value: string): SafeHtml {
    return DOMPurify.sanitize(value);
  }

  ngOnChanges(changes: { [propKey: string]: SimpleChange }): void {
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
