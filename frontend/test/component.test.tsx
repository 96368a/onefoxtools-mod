import { cleanup, fireEvent, render } from 'solid-testing-library'
import { afterEach, describe, expect, it } from 'vitest'
import Counter from '../src/components/Counter'

describe('Counter', () => {
  afterEach(cleanup)

  it('should render', () => {
    const { queryByText } = render(() => <Counter initial={10}/>)
    expect(queryByText('10')).toBeDefined()
  })

  it('should be interactive', () => {
    const { queryByText } = render(() => <Counter initial={0}/>)
    expect(queryByText('0')).toBeDefined()

    fireEvent.click(queryByText('+') as HTMLElement)

    expect(queryByText('1')).toBeDefined()
  })
})
