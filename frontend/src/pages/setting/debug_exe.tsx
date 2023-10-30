import { TestCmdExec } from 'wailsjs/go/main/CONFIG'
import type { main } from 'wailsjs/go/models'
import DataStore from '~/store/data'

export default function () {
  const { configs } = DataStore
  const tools = createMemo(() => {
    return configs.reduce((tools, c) => tools.concat(c.config), [] as main.Config[])
  })
  const [status, setStatus] = createStore<Array<number>>(Array.from({ length: tools().length }, () => 0))

  function TestExe(i: number) {
    TestCmdExec(tools()[i]).then(() => {
      setStatus(i, 1)
    }).catch(() => {
      setStatus(i, -1)
    })
  }
  return (
        <div class="h-full w-full bg-base-300 pt-4 rounded-box">
            <div class="py-2">
                工具启动调试界面
            </div>
            <div>
                <button class='btn'>开始测试</button>
                <button class='btn btn-warning'>停止测试</button>

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
                                        <td>
                                            <Switch fallback={<div class="icon-btn i-carbon-face-neutral mx-auto h-6 w-6" onclick={() => TestExe(i())} />}>
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
        </div>
  )
}
