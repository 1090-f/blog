import { marked } from 'marked'
import katex from 'katex'

marked.setOptions({
  gfm: true,
  breaks: true
})

// 渲染 Markdown，并在交给 marked 前保护和解析 KaTeX 公式。
export function renderMarkdown(source = '') {
  const formulas = []
  let markdown = String(source)

  markdown = markdown.replace(/\$\$([\s\S]+?)\$\$/g, (_, expression) => {
    const index = formulas.push({ expression: expression.trim(), displayMode: true }) - 1
    return `\n\nKATEXBLOCKTOKEN${index}END\n\n`
  })

  markdown = markdown.replace(/\$([^\n$]+?)\$/g, (_, expression) => {
    const index = formulas.push({ expression: expression.trim(), displayMode: false }) - 1
    return `KATEXINLINETOKEN${index}END`
  })

  let html = marked(markdown)
  formulas.forEach(({ expression, displayMode }, index) => {
    const token = displayMode ? `KATEXBLOCKTOKEN${index}END` : `KATEXINLINETOKEN${index}END`
    const rendered = katex.renderToString(expression, {
      displayMode,
      throwOnError: false,
      strict: false
    })
    html = displayMode
      ? html.replace(`<p>${token}</p>`, rendered)
      : html.replaceAll(token, rendered)
  })

  return html
}
