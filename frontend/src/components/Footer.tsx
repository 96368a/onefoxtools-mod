import { WindowSetDarkTheme, WindowSetLightTheme } from 'wailsjs/runtime/runtime'

export default function Footer() {
  const [isDark, setIsDark] = createSignal(false)
  const navigate = useNavigate()
  function toggle() {
    // toggleDark()
    setIsDark(document.documentElement.getAttribute('data-theme') === 'dark')
    if (isDark()) {
      WindowSetLightTheme()
      document.documentElement.setAttribute('data-theme', 'light')
    }
    else {
      WindowSetDarkTheme()
      document.documentElement.setAttribute('data-theme', 'dark')
    }
  }
  return (
    <nav class="mt-6 inline-flex gap-2 text-xl">
      <button
        class="icon-btn i-carbon-home"
        title="主页"
        onclick={() => navigate('/')}
      ></button>
      <button class="icon-btn !outline-none" onClick={toggle}>
        {isDark() ? <div class="i-carbon-moon" /> : <div class="i-carbon-sun" />}
      </button>
      <a
        class="icon-btn i-carbon-logo-github"
        rel="noreferrer"
        href="https://github.com/96368a"
        target="_blank"
        title="GitHub"
      />

    </nav>

  )
}
