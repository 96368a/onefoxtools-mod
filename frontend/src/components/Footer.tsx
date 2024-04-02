import { BrowserOpenURL, WindowSetDarkTheme, WindowSetLightTheme } from '../../wailsjs/runtime/runtime'

export default function Footer() {
  const [isDark, setIsDark] = createSignal(false)
  const navigate = useNavigate()
  function toggle() {
    // toggleDark()
    setIsDark(document.documentElement.getAttribute('data-theme') === 'dark')
    if (isDark()) {
      document.documentElement.setAttribute('data-theme', 'light')
      WindowSetLightTheme()
    }
    else {
      document.documentElement.setAttribute('data-theme', 'dark')
      WindowSetDarkTheme()
    }
  }
  return (
    <nav class="inline-flex justify-center gap-2 py-6 text-xl">
      <button
        class="icon-btn i-carbon-home"
        title="主页"
        onclick={() => navigate('/')}
      ></button>
      <button
        title='切换主题'
        class="icon-btn !outline-none"
        onClick={toggle}>
        {isDark() ? <div class="i-carbon-moon" /> : <div class="i-carbon-sun" />}
      </button>
      <button
        class="icon-btn i-carbon-logo-github"
        title="项目地址"
        onclick={() => BrowserOpenURL('https://github.com/96368a/OnefoxTools-Mod')}
      ></button>
    </nav>

  )
}
