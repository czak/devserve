new EventSource('/events').addEventListener('change', e => {
  const path = '/' + e.data

  for (const link of document.getElementsByTagName("link")) {
    const url = new URL(link.href)

    if (url.host === location.host && url.pathname === path) {
      const next = link.cloneNode()
      next.href = path + '?' + Math.random().toString(36).slice(2)
      next.onload = () => link.remove()
      link.parentNode.insertBefore(next, link.nextSibling)
      return
    }
  }

  location.reload()
})
