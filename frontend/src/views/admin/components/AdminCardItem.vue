<template>
  <li class="admin-card relative overflow-hidden active:scale-[0.98] transition-all duration-200">
    <slot name="background"></slot>
    <slot name="overlay"></slot>
    <div class="relative z-10 block p-5">
      <div>
        <slot name="header"></slot>
      </div>
      <div class="mt-4">
        <slot name="content"></slot>
      </div>
      <div class="mt-8">
        <slot name="footer"></slot>
      </div>
    </div>
  </li>
</template>

<script setup lang="ts"></script>
