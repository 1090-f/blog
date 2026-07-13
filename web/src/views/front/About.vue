<template>
  <div class="about-page">
    <section class="about-intro">
      <h1 class="about-title">{{ copy.title }}</h1>
      <p class="about-lead">
        {{ copy.greeting }} <strong>{{ authorName }}</strong>，{{ copy.lead }}
      </p>

      <ul class="about-list">
        <li>
          <span class="about-bullet" aria-hidden="true"></span>
          <span>{{ copy.todo }}</span>
          <a href="https://bangumi.tv/" target="_blank" rel="noreferrer">bangumi 番组计划</a>
        </li>
        <li>
          <span class="about-bullet" aria-hidden="true"></span>
          <span>{{ copy.tools }}</span>
        </li>
      </ul>

      <div class="tool-grid" aria-label="常用工具">
        <div v-for="tool in tools" :key="tool.name" class="tool-item" :class="`tool-${tool.theme}`" :title="tool.name">
          <span>{{ tool.mark }}</span>
        </div>
      </div>
    </section>

    <section class="about-section">
      <h2>{{ copy.siteTitle }}</h2>
      <p class="about-section-lead">
        {{ copy.siteLead }} <strong>Gin Blog</strong>。
      </p>

      <article class="project-card">
        <div class="project-heading">
          <div class="project-identity">
            <img src="/author-profile.png?v=2" alt="Gin Blog avatar" class="project-avatar" />
            <span>Gin Blog / <strong>Firefly</strong></span>
          </div>
          <a class="project-github" href="https://github.com/" target="_blank" rel="noreferrer" aria-label="GitHub">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M12 3a9 9 0 0 0-2.8 17.55c.45.08.62-.2.62-.44v-1.55c-2.53.55-3.06-1.07-3.06-1.07-.4-1.02-.98-1.29-.98-1.29-.8-.55.06-.54.06-.54.89.07 1.36.92 1.36.92.79 1.36 2.07.97 2.58.74.08-.57.31-.97.56-1.19-2.02-.23-4.15-1.01-4.15-4.5 0-.99.35-1.8.92-2.43-.09-.22-.4-1.15.09-2.4 0 0 .75-.24 2.47.93A8.6 8.6 0 0 1 12 7.4c.76 0 1.52.1 2.23.33 1.72-1.17 2.47-.93 2.47-.93.49 1.25.18 2.18.09 2.4.57.63.92 1.44.92 2.43 0 3.5-2.13 4.26-4.16 4.49.32.28.6.82.6 1.65v2.45c0 .24.16.52.62.44A9 9 0 0 0 12 3Z" />
            </svg>
          </a>
        </div>
        <p class="project-description">{{ copy.projectDescription }}</p>
        <div class="project-meta">
          <span>☆ 1.6K</span>
          <span>⌘ 1.1K</span>
          <span>⚖ MIT</span>
        </div>
      </article>
    </section>

    <section class="about-section contact-section">
      <h2>{{ copy.contactTitle }}</h2>
      <p class="about-section-lead">{{ copy.contactLead }}</p>
      <div class="contact-links">
        <a v-for="link in contactLinks" :key="link.label" :href="link.href" :target="link.external ? '_blank' : undefined" :rel="link.external ? 'noreferrer' : undefined">
          {{ link.label }}
        </a>
      </div>
    </section>
  </div>
</template>

<script setup>
import { useUserStore } from '../../stores/user'

const userStore = useUserStore()
const authorName = userStore.user?.nickname || 'Jason Shane'

const copy = {
  title: '\u5173\u4e8e\u6211 / About Me',
  greeting: '\u4f60\u597d\uff01\u6211\u662f',
  lead: '\u4e00\u4e2a\u4e8c\u6b21\u5143\u6280\u672f\u5b85\uff0c\u559c\u6b22\u63a2\u7d22\u65b0\u7684\u6280\u672f\u548c\u586b\u5145\u81ea\u5df1\u7684\u672a\u77e5\u57df\uff0c\u6311\u6218\u4e0d\u53ef\u77e5\u57df\u3002',
  todo: '\u8ffd\u756a\u6e05\u5355\uff1a',
  tools: '\u5e38\u7528\u5de5\u5177\uff1a',
  siteTitle: '\u5173\u4e8e\u672c\u7ad9',
  siteLead: '\u672c\u7ad9\u4f7f\u7528',
  projectDescription: '\u8bb0\u5f55\u6280\u672f\u3001\u751f\u6d3b\u548c\u6bcf\u4e00\u6b21\u601d\u8003\u3002\u8fd9\u91cc\u662f\u5c5e\u4e8e\u6211\u7684\u5b89\u9759\u89d2\u843d\u3002',
  contactTitle: '\u8054\u7cfb\u65b9\u5f0f',
  contactLead: '\u5982\u679c\u4f60\u60f3\u548c\u6211\u4ea4\u6d41\u6280\u672f\u95ee\u9898\u3001\u5206\u4eab\u6709\u8da3\u7684\u60f3\u6cd5\uff0c\u6b22\u8fce\u901a\u8fc7\u4ee5\u4e0b\u65b9\u5f0f\u8054\u7cfb\u6211\u3002'
}

const tools = [
  { name: 'Python', mark: 'Py', theme: 'python' },
  { name: 'Markdown', mark: 'M↓', theme: 'markdown' },
  { name: 'Git', mark: '◆', theme: 'git' },
  { name: 'GitHub', mark: '●', theme: 'github' },
  { name: 'Cloudflare', mark: '☁', theme: 'cloudflare' },
  { name: 'Astro', mark: '✦', theme: 'astro' },
  { name: 'VS Code', mark: '⌁', theme: 'vscode' },
  { name: 'Blender', mark: '◎', theme: 'blender' },
  { name: 'Notion', mark: 'N', theme: 'notion' },
  { name: 'Obsidian', mark: '◈', theme: 'obsidian' }
]

const contactLinks = [
  { label: 'GitHub', href: 'https://github.com/', external: true },
  { label: 'Bilibili', href: 'https://www.bilibili.com/', external: true },
  { label: 'hello@ginblog.dev', href: 'mailto:hello@ginblog.dev', external: false }
]
</script>

<style scoped>
.about-page {
  position: relative;
  max-width: 1180px;
  margin: 0 auto;
  padding: 20px 0 72px;
}

.about-page::before,
.about-page::after {
  content: '';
  position: absolute;
  width: 18px;
  height: 34px;
  border-radius: 70% 30% 70% 30%;
  background: rgba(244, 180, 194, 0.5);
  transform: rotate(28deg);
  filter: blur(0.5px);
  pointer-events: none;
}

.about-page::before {
  top: -30px;
  right: 62%;
}

.about-page::after {
  top: 300px;
  right: 36%;
  transform: rotate(48deg) scale(0.7);
}

.about-intro {
  padding: 0 0 34px;
}

.about-title,
.about-section h2 {
  color: var(--text-primary);
  font-size: clamp(36px, 5vw, 56px);
  line-height: 1.1;
  letter-spacing: -0.04em;
}

.about-title {
  margin-bottom: 42px;
}

.about-lead,
.about-section-lead {
  max-width: 1100px;
  color: #c7cbd1;
  font-size: clamp(18px, 2vw, 25px);
  line-height: 1.8;
}

.about-lead strong,
.about-section-lead strong {
  color: #fff;
}

.about-list {
  display: grid;
  gap: 14px;
  margin: 30px 0 22px;
  padding: 0;
  list-style: none;
  color: #c7cbd1;
  font-size: 21px;
}

.about-list li {
  display: flex;
  align-items: center;
  gap: 16px;
}

.about-bullet {
  width: 8px;
  height: 8px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: var(--accent);
}

.about-list a,
.contact-links a {
  color: var(--accent);
}

.about-list a:hover,
.contact-links a:hover {
  color: #fff;
}

.tool-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
  max-width: 850px;
  margin: 46px auto 0;
}

.tool-item {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  border-radius: 16px;
  background: #f3f3f0;
  color: #20242a;
  font-size: 23px;
  font-weight: 800;
  letter-spacing: -0.08em;
  box-shadow: inset 0 -2px 0 rgba(0, 0, 0, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.tool-item:hover {
  transform: translateY(-5px) rotate(-2deg);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.25);
}

.tool-python { color: #e3b52d; text-shadow: 2px 2px #3979a8; }
.tool-markdown { color: #161616; font-size: 21px; }
.tool-git { color: #f04438; font-size: 35px; }
.tool-github { color: #191919; font-size: 40px; }
.tool-cloudflare { color: #f28b20; font-size: 38px; }
.tool-astro { color: #4e83b7; font-size: 38px; }
.tool-vscode { color: #1478bd; font-size: 42px; }
.tool-blender { color: #e87900; font-size: 42px; }
.tool-notion { color: #171717; border: 2px solid #1a1a1a; }
.tool-obsidian { color: #8054de; font-size: 40px; }

.about-section {
  margin-top: 70px;
}

.about-section h2 {
  margin-bottom: 34px;
  font-size: clamp(34px, 4vw, 48px);
}

.project-card {
  margin-top: 30px;
  padding: 32px 36px;
  border-radius: 24px;
  background: rgba(3, 14, 11, 0.72);
  border: 1px solid rgba(255, 255, 255, 0.025);
}

.project-heading,
.project-identity,
.project-meta,
.contact-links {
  display: flex;
  align-items: center;
}

.project-heading {
  justify-content: space-between;
  gap: 20px;
}

.project-identity {
  gap: 14px;
  color: #e7eaeb;
  font-size: 26px;
}

.project-identity strong {
  color: #fff;
}

.project-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  object-position: center 13%;
}

.project-github {
  color: #e7eaeb;
}

.project-github:hover {
  color: var(--accent);
}

.project-github svg {
  width: 34px;
  height: 34px;
  fill: currentColor;
}

.project-description {
  margin: 20px 0 18px;
  color: #b9c0c4;
  font-size: 18px;
  line-height: 1.8;
}

.project-meta {
  gap: 28px;
  color: #c7cdd0;
  font-size: 16px;
}

.contact-section {
  padding-bottom: 20px;
}

.contact-links {
  flex-wrap: wrap;
  gap: 24px;
  margin-top: 26px;
  font-size: 18px;
}

@media (max-width: 768px) {
  .about-page {
    padding-top: 10px;
  }

  .about-title {
    margin-bottom: 26px;
  }

  .about-list {
    font-size: 17px;
  }

  .tool-item {
    width: 58px;
    height: 58px;
    border-radius: 13px;
    font-size: 18px;
  }

  .project-card {
    padding: 24px 20px;
    border-radius: 18px;
  }

  .project-identity {
    font-size: 20px;
  }

  .project-description {
    font-size: 15px;
  }

  .project-meta {
    flex-wrap: wrap;
    gap: 14px;
    font-size: 14px;
  }
}
</style>
