import { Injectable, OnDestroy, Renderer2, RendererFactory2 } from '@angular/core';
import { fromEventPattern, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class StyleManagerService implements OnDestroy {
  private _destroy$ = new Subject();
  private pointer = 0;
  private keyArray: string[] = [
    'ArrowUp',
    'ArrowUp',
    'ArrowDown',
    'ArrowDown',
    'ArrowLeft',
    'ArrowRight',
    'ArrowLeft',
    'ArrowRight',
    'b',
    'a',
    'Enter',
  ];

  constructor(private rendererFactory2: RendererFactory2) {
    const renderer = this.rendererFactory2.createRenderer(null, null);
    this.createOnClickObservable(renderer);
  }

  private createOnClickObservable(renderer: Renderer2): void {
    let removeClickEventListener: () => void;
    const createClickEventListener = (handler: (e: KeyboardEvent) => boolean | void): void => {
      removeClickEventListener = renderer.listen('document', 'keydown', handler);
    };

    fromEventPattern<Event>(createClickEventListener, () => removeClickEventListener())
      .pipe(takeUntil(this._destroy$))
      .subscribe((event: Event) => {
        this.handleKeyboardEvent(event as KeyboardEvent);
      });
  }

  private handleKeyboardEvent(event: KeyboardEvent): void {
    if (this.keyArray[this.pointer] === event.key) {
      ++this.pointer;
      if (this.pointer === this.keyArray.length) {
        this.setDarkMode();
        this.pointer = 0;
      }
    } else {
      this.pointer = 0;
    }
  }

  private setDarkMode(): void {
    this.setStyle('theme', 'assets/default-branding/dark-theme.css');
    // TODO: change to 'branding' instead of 'default-branding'
  }

  public setStyle(key: string, href: string): void {
    this.getLinkElementForKey(key).setAttribute('href', href);
  }

  private getLinkElementForKey(key: string): HTMLLinkElement {
    return this.getExistingLinkElementByKey(key) || this.createLinkElementWithKey(key);
  }

  private getExistingLinkElementByKey(key: string): HTMLLinkElement | null {
    return document.head.querySelector(`link[rel="stylesheet"].${this.getClassNameForKey(key)}`);
  }

  private createLinkElementWithKey(key: string): HTMLLinkElement {
    const linkEl = document.createElement('link');
    linkEl.setAttribute('rel', 'stylesheet');
    linkEl.classList.add(this.getClassNameForKey(key));
    document.body.appendChild(linkEl);
    return linkEl;
  }

  private getClassNameForKey(key: string): string {
    return `app-${key}`;
  }

  public ngOnDestroy(): void {
    this._destroy$.next();
    this._destroy$.complete();
  }
}
