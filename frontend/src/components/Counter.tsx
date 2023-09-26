export default function Counter({ initial }: { initial: number }) {
  const [count, setCount] = createSignal(initial ?? 0)
  function increment() {
    setCount(count() + 1)
  }

  function decrement() {
    setCount(count() - 1)
  }

  return (
    <div>
      {count()}
      <button class="inc" onClick={() => increment()}>
        +
      </button>
      <button class="dec" onClick={() => decrement()}>
        -
      </button>
    </div>
  )
}
