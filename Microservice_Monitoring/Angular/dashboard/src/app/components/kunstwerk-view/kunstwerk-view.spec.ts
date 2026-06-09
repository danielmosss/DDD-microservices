import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KunstwerkViewComponent } from './kunstwerk-view';

describe('Kusntwerkview', () => {
  let component: KunstwerkViewComponent;
  let fixture: ComponentFixture<KunstwerkViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KunstwerkViewComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(KunstwerkViewComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
