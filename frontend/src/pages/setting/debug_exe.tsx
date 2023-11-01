import { TestCmdExec } from 'wailsjs/go/main/CONFIG'
import type { main } from 'wailsjs/go/models'
import DataStore from '~/store/data'

export default function () {
  const { configs } = DataStore
  const tools = createMemo(() => {
    return configs.reduce((tools, c) => tools.concat(c.config), [] as main.Config[])
  })
  const [status, setStatus] = createStore<Array<number>>(Array.from({ length: tools().length }, () => 0))
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-nocheck
  let dialog: HTMLDialogElement
  function TestExe(i: number) {
    TestCmdExec(tools()[i]).then(() => {
      setStatus(i, 1)
    }).catch(() => {
      setStatus(i, -1)
    })
  }
  const [flag, setFlag] = createSignal(false)

  async function StartTest() {
    setFlag(false)
    const configs = Array.from(tools())
    TestAll(configs, 0)
  }

  async function TestAll(configs: main.Config[], i: number) {
    try {
      if (flag())
        return
      if (configs.length === 0)
        return
      await TestCmdExec(configs.shift()!)
      setStatus(i, 1)
    }
    catch (e) {
      setStatus(i, -1)
    }
    finally {
      TestAll(configs, i + 1)
    }
  }

  return (
        <div class="h-full w-full bg-base-300 pt-4 rounded-box">
            <div class="py-2">
                工具启动调试界面
            </div>
            <div>
                <button class='mr-10 btn' onclick={() => dialog.showModal()}>开始测试</button>
                <button class='btn btn-warning' onclick={() => setFlag(true)}>停止测试</button>
            </div>
            <div class="flex flex-col gap-2 p-2">
                <div class="overflow-x-auto">
                    <table class="w-full table">
                        <thead>
                            <tr>
                                <th></th>
                                <th>分类</th>
                                <th>名称</th>
                                <th>状态</th>
                            </tr>
                        </thead>
                        <tbody>
                            <For each={tools()}>{
                                (c, i) => (
                                    <tr>
                                        <th>{i() + 1}</th>
                                        <td>111</td>
                                        <td>{c.name}</td>
                                        <td onclick={() => TestExe(i())}>
                                            <Switch fallback={<div class="icon-btn i-carbon-face-neutral mx-auto h-6 w-6" />}>
                                                <Match when={status[i()] === 1}>
                                                    <div class="icon-btn i-carbon-face-satisfied mx-auto h-6 w-6 color-green" />
                                                </Match>
                                                <Match when={status[i()] === -1}>
                                                    <div class="icon-btn i-carbon-face-dissatisfied mx-auto h-6 w-6 color-red" />
                                                </Match>
                                            </Switch>
                                        </td>
                                    </tr>
                                )
                            }</For>
                        </tbody>
                    </table>
                </div>
            </div>
            <dialog ref={dialog} id="dialog_warn" class="modal">
                <div class="modal-box">
                    <h3 class="text-lg font-bold">⚠️警告⚠️</h3>
                    <p class="py-4">测试会把所有工具全都启动，可能会导致电脑卡死！！！</p>
                    <p class="py-4">部分工具测试完毕之后会自动关闭，但是大量工具仍然需要手动关闭</p>
                    <p class="py-4">测试过程中你可以随时停止</p>
                    <div class="justify-center modal-action">
                        <form method="dialog" class='flex gap-4'>
                            <button class="w-30 btn btn-warning" onclick={StartTest}>⚡确定⚡</button>
                            <button class="w-30 btn">取消</button>
                        </form>
                    </div>
                </div>
            </dialog>
        </div>
  )
}
