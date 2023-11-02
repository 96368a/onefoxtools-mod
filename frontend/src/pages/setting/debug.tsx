import { WindowSetTitle } from 'wailsjs/runtime/runtime'

export default function () {
  function setTitle() {
    WindowSetTitle('test')
  }
  return (
          <div class="h-full w-full bg-base-300 pt-4 rounded-box">
              set
              <button class="btn" onclick={setTitle}>设置标题</button>
          </div>
  )
}
