/* @refresh reload */
import { render } from 'solid-js/web'

import App from './App'
import '@unocss/reset/tailwind.css'
import 'uno.css'
import './styles/main.css'

render(() => <App />, document.getElementById('root') as HTMLElement)
