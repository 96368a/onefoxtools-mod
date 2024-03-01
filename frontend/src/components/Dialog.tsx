import type { JSX } from 'solid-js'

const [props, setProps] = createStore<DialogProps>({} as DialogProps)
const [isShow, setShow] = createSignal(false)

let dialogsRef: HTMLDialogElement
export default function Dialog() {
  // 根据show控制显示
  createEffect(() => {
    if (isShow())
      dialogsRef?.showModal()
    else
      dialogsRef?.close()
  })
  return <div>
    <Portal mount={document.body}>
      {/* <button class="btn" onclick={() => show('hello', 'world')}>open modal</button> */}
      <dialog ref={d => dialogsRef = d} class="modal" onClose={() => setShow(false)}>
        <div class="modal-box">
          <form method="dialog">
            <button class="absolute right-2 top-2 btn btn-circle btn-ghost btn-sm">✕</button>
          </form>
          <h3 class="text-lg font-bold">{props.title}</h3>
          <p class="break-all py-4">{props.msg}</p>
        </div>
        <form method="dialog" class="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </Portal>
  </div>
}

// 设置对话框属性并显示
export function show(title: string, msg = '', onClose?: () => void, children?: JSX.Element) {
  setProps({ title, msg, onClose, children })
  setShow(true)
}

export function close() {
  setShow(false)
}

interface DialogProps {
  title: string
  msg: string
  onClose?: () => void
  children?: JSX.Element
}
