import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { environment } from 'src/environments/environment';
import { AppComponent } from './app.component';

describe('AppComponent', () => {
  let fixture: ComponentFixture<AppComponent>;
  let app: AppComponent;
  let compiled: HTMLElement;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        RouterTestingModule
      ],
      declarations: [
        AppComponent
      ],
    }).compileComponents();
  });

  beforeEach(async () => {
    fixture = TestBed.createComponent(AppComponent);
    app = fixture.componentInstance;
    fixture.detectChanges();
    compiled = fixture.nativeElement as HTMLElement;
  });

  it('should create the app', () => {    
    expect(app).toBeTruthy();
  });

  it(`should have as title 'demo-app'`, () => {
    expect(app.title).toEqual('demo-app');
  });

  it('should render title', () => {
    const title = compiled.querySelector('[data-test="title"]')
    
    expect(title?.textContent).toEqual(environment.title);
  });

  it('should render version', () => {
    const version = compiled.querySelector('[data-test="version"]');
    const versionTitle = compiled.querySelector('[data-test="version-title"]');

    expect(version?.textContent).toEqual(environment.version);
    expect(versionTitle?.textContent).toEqual('Version: ' + environment.version);
  });

  it('should render production text', () => {
    const production = compiled.querySelector('[data-test="production"]');
    const productionTitle = compiled.querySelector('[data-test="production-title"]');

    expect(production?.textContent).toEqual(' not ');
    expect(productionTitle?.textContent).toEqual('Running not in production');
  });
});
