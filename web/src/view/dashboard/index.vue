<template>
  <div
    class="h-full gva-container2 overflow-auto bg-white text-black dark:bg-slate-800 dark:text-white"
  >
    <div class="p-4 lg:p-6">
      <div class="grid grid-cols-1 gap-4">
        <gva-card title="快捷功能" show-action>
          <gva-quick-link />
        </gva-card>
      </div>
    </div>
  </div>
</template>

<script setup>
  import {
    GvaQuickLink,
    GvaCard
  } from './components'

  defineOptions({
    name: 'Dashboard'
  })
</script>

<style lang="scss" scoped></style>
