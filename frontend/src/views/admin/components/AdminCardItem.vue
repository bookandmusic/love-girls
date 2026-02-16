<template>
  <li class="admin-card relative overflow-hidden">
    <slot name="background"></slot>
    <slot name="overlay"></slot>
    <div class="relative z-10 block p-4">
      <div>
        <slot name="header"></slot>
      </div>
      <div class="mt-3">
        <slot name="content"></slot>
      </div>
      <div class="mt-6">
        <slot name="footer"></slot>
      </div>
    </div>
  </li>
</template>

<script setup lang="ts"></script>
