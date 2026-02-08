<template>
  <div class="relative w-full min-h-screen">
    <div
      ref="container"
      class="absolute inset-0 overflow-hidden z-0 bg-gradient-to-b from-white to-[var(--primary-color)]"
    ></div>
    <div class="relative z-10 w-full min-h-screen">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import p5 from 'p5'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

/* ======================
   Props 定义
====================== */
interface Props {
  isPlaying?: boolean
  isDarkMode?: boolean
  waveColor: string
}

const props = withDefaults(defineProps<Props>(), {
  isPlaying: true,
  isDarkMode: false,
})

const container = ref<HTMLElement | null>(null)

// 扩展p5类型以包括我们自定义的方法
interface Extendedp5 extends p5 {
  updateProps?: (props: Props) => void
}

let p5Instance: Extendedp5 | null = null

/* ======================
   p5 Sketch
====================== */
const sketch = (p: Extendedp5) => {
  let mountains: Mountain[] = []
  let bgColor = '#e6e6e6'
  let isDarkMode = false
  let waveColor = props.waveColor

  p.setup = () => {
    p.createCanvas(p.windowWidth, p.windowHeight)
    mountains = []
    growMountains(p, mountains, waveColor)
    p.background(bgColor)
    mountains.forEach(m => m.display(p))
  }

  p.draw = () => {
    p.background(bgColor)
    mountains.forEach(m => m.display(p))
  }

  p.windowResized = () => {
    p.resizeCanvas(p.windowWidth, p.windowHeight)
  }

  /* ===========
       对应 React 的 updateWithProps
    =========== */
  p.updateProps = (newProps: Props) => {
    p.frameRate(newProps.isPlaying ? 30 : 0)

    bgColor = newProps.isDarkMode ? '#323232' : '#e6e6e6'

    if (isDarkMode !== (newProps.isDarkMode ?? false) || waveColor !== newProps.waveColor) {
      isDarkMode = newProps.isDarkMode ?? false
      waveColor = newProps.waveColor
      p.setup()
    }
  }
}

/* ======================
   生命周期
====================== */
onMounted(() => {
  if (container.value && props.isPlaying) {
    p5Instance = new p5(sketch, container.value) as Extendedp5
  }
})

onBeforeUnmount(() => {
  p5Instance?.remove()
  p5Instance = null
})

/* ======================
   Props 监听
====================== */
watch(
  () => ({ ...props }),
  newProps => {
    if (p5Instance && p5Instance.updateProps) {
      p5Instance.updateProps(newProps)
    }
  },
  { deep: true }
)

/* ======================
   Mountain 类 & 工具函数
====================== */
class Mountain {
  c: p5.Color
  y: number
  offset: number
  t: number

  constructor(color: p5.Color, y: number, p: p5) {
    this.c = color
    this.y = y
    this.offset = p.random(100, 200)
    this.t = 0
  }

  display(p: p5) {
    let xoff = 0

    p.noStroke()
    p.fill(this.c)
    p.noiseDetail(1.7, 1.3)

    p.beginShape()
    for (let x = 0; x <= p.width + 25; x += 25) {
      const yoff = p.map(p.noise(xoff + this.offset, this.t + this.offset), 0, 1, 0, 200)
      const y = this.y - yoff
      p.vertex(x, y)
      xoff += 0.08
    }

    p.vertex(p.width + 100, p.height)
    p.vertex(0, p.height)
    p.endShape(p.CLOSE)

    this.t += 0.005
  }
}

function growMountains(p: p5, mountains: Mountain[], hexColor: string) {
  const c = p.color(hexColor)

  new Array(5).fill(1).forEach((_, i) => {
    const alpha = 255 - 50 * i
    c.setAlpha(alpha)
    const h = p.height - 50 * i
    mountains.push(new Mountain(c, h, p))
  })
}
</script>
