export default function Footer() {
  const { isDark, toggleDark } = useDark()
  function toggle() {
    toggleDark()
    document.documentElement.setAttribute('data-theme', isDark() ? 'dark' : 'light')
  }
  return (
    <nav class="text-xl mt-6 inline-flex gap-2">
      <button class="icon-btn !outline-none" onClick={ toggle }>
        {isDark() ? <div class="i-carbon-moon" /> : <div class="i-carbon-sun" />}
      </button>
      <button data-toggle-theme="dark,light" data-act-class="ACTIVECLASS">111</button>
      <button data-set-theme="dark" data-act-class="ACTIVECLASS">dark</button>
      <a
        class="icon-btn i-carbon-logo-github"
        rel="noreferrer"
        href="https://github.com/nanakura/vitesse-lite-solidjs"
        target="_blank"
        title="GitHub"
      />
    </nav>

  )
}
