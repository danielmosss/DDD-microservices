import { Routes } from '@angular/router';
import { KunstwerkViewComponent } from './components/kunstwerk-view/kunstwerk-view';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'dashboard',
    pathMatch: 'full'
  },
  {
    path: 'dashboard',
    component: KunstwerkViewComponent
  },
  // Optioneel: een catch-all voor als je een onbekende URL intypt
  {
    path: '**',
    redirectTo: 'dashboard'
  }
];
