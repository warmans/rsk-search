import { Directive, ElementRef, Input } from '@angular/core';

@Directive({
  selector: '[focused]'
})
export class FocusedDirective {
  @Input()
  set focused(value: boolean) {
    if (value) {
      this.elementRef.nativeElement.scrollIntoViewIfNeeded();
    }
  }

  constructor(private elementRef: ElementRef) {
  }
}
