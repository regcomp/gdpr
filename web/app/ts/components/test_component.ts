const boop = () => {
  return {
    chirp() {
      console.log('CACAW');
    }
  }
}

document.addEventListener('alpine:init', () => {
  window.Alpine.data('boop', boop)
})
