import { Routes } from '@angular/router';
import { KunstwerkViewComponent } from './components/kunstwerk-view/kunstwerk-view';
import { Kunstwerkkiezer } from './components/kunstwerkkiezer/kunstwerkkiezer';

export const routes: Routes = [
  {
    path: '',
    component: Kunstwerkkiezer
  },
  {
    path: 'kunstwerk/:kunstwerkid',
    component: KunstwerkViewComponent
  },
  // Optioneel: een catch-all voor als je een onbekende URL intypt
  {
    path: '**',
    redirectTo: 'dashboard'
  }
];
