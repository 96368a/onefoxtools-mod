export default function Footer() {
  const { isDark, toggleDark } = useDark()
  return (
    <nav class="text-xl mt-6 inline-flex gap-2">
      <button class="icon-btn !outline-none" onClick={() => toggleDark()}>
        {isDark() ? <div class="i-carbon-moon" /> : <div class="i-carbon-sun" />}
      </button>

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
