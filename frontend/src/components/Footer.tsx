export default function Footer() {
  const [isDark, setIsDark] = createSignal(false)
  function toggle() {
    // toggleDark()
    setIsDark(document.documentElement.getAttribute('data-theme') === 'dark')
    document.documentElement.setAttribute('data-theme', isDark() ? 'light' : 'dark')
  }
  return (
    <nav class="text-xl mt-6 inline-flex gap-2">
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
