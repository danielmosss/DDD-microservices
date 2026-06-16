import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Kunstwerkkiezer } from './kunstwerkkiezer';

describe('Kunstwerkkiezer', () => {
  let component: Kunstwerkkiezer;
  let fixture: ComponentFixture<Kunstwerkkiezer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Kunstwerkkiezer],
    }).compileComponents();

    fixture = TestBed.createComponent(Kunstwerkkiezer);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
