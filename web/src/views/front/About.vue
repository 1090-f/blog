<template>
  <div class="about-page">
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
      </article>

      <div class="site-details">
        <article v-for="item in siteHighlights" :key="item.title" class="site-detail-card">
          <span class="site-detail-icon" aria-hidden="true">{{ item.icon }}</span>
          <h3>{{ item.title }}</h3>
          <p>{{ item.description }}</p>
        </article>
      </div>
    </section>

    <section class="about-section contact-section">
      <h2>{{ copy.contactTitle }}</h2>
      <p class="about-section-lead">{{ copy.contactLead }}</p>
      <div class="contact-links">
        <a v-for="link in contactLinks" :key="link.label" :href="link.href" :target="link.external ? '_blank' : undefined" :rel="link.external ? 'noreferrer' : undefined">
          {{ link.label }}
        </a>
      </div>
      <p class="contact-note">添加好友时请注明“博客来访”，方便确认你的留言。</p>
    </section>
  </div>
</template>

<script setup>
const copy = {
  siteTitle: '\u5173\u4e8e\u672c\u7ad9',
  siteLead: '\u672c\u7ad9\u4f7f\u7528',
  projectDescription: '\u8bb0\u5f55\u6280\u672f\u3001\u751f\u6d3b\u548c\u6bcf\u4e00\u6b21\u601d\u8003\u3002\u8fd9\u91cc\u662f\u5c5e\u4e8e\u6211\u7684\u5b89\u9759\u89d2\u843d\u3002',
  contactTitle: '\u8054\u7cfb\u65b9\u5f0f',
  contactLead: '\u5982\u679c\u4f60\u60f3\u548c\u6211\u4ea4\u6d41\u6280\u672f\u95ee\u9898\u3001\u5206\u4eab\u6709\u8da3\u7684\u60f3\u6cd5\uff0c\u6b22\u8fce\u901a\u8fc7\u4ee5\u4e0b\u65b9\u5f0f\u8054\u7cfb\u6211\u3002'
}

const contactLinks = [
  {
    label: 'QQ：1438318243',
    href: 'tencent://AddContact/?fromId=45&fromSubId=1&subcmd=all&uin=1438318243&website=www.oicqzone.com',
    external: false
  }
]

const siteHighlights = [
  {
    icon: '✦',
    title: '内容方向',
    description: '记录开发实践、工具使用与日常思考，也会整理值得回看的技术笔记。'
  },
  {
    icon: '◌',
    title: '写作原则',
    description: '尽量写清背景、过程与结论；持续补充和修订，让每篇文章更有参考价值。'
  },
  {
    icon: '↗',
    title: '交流互动',
    description: '欢迎在文章下留言讨论，也欢迎通过 QQ 分享建议、问题或有趣的想法。'
  }
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

.about-section h2 {
  color: var(--text-primary);
  font-size: clamp(36px, 5vw, 56px);
  line-height: 1.1;
  letter-spacing: -0.04em;
}

.about-section-lead {
  max-width: 1100px;
  color: #c7cbd1;
  font-size: clamp(18px, 2vw, 25px);
  line-height: 1.8;
}

.about-section-lead strong {
  color: #fff;
}

.contact-links a {
  color: var(--accent);
}

.contact-links a:hover {
  color: #fff;
}

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

.site-details {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.site-detail-card {
  min-height: 190px;
  padding: 24px;
  border: 1px solid var(--border);
  border-radius: 18px;
  background: rgba(17, 28, 43, 0.55);
}

.site-detail-icon {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  margin-bottom: 18px;
  border-radius: 10px;
  background: var(--accent-dim);
  color: var(--accent);
  font-size: 19px;
  font-weight: 700;
}

.site-detail-card h3 {
  margin-bottom: 10px;
  color: var(--text-primary);
  font-size: 18px;
}

.site-detail-card p {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.75;
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

.contact-note {
  margin-top: 16px;
  color: var(--text-muted);
  font-size: 14px;
}

@media (max-width: 768px) {
  .about-page {
    padding-top: 10px;
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

  .site-details {
    grid-template-columns: 1fr;
  }

  .site-detail-card {
    min-height: 0;
    padding: 20px;
  }

}
</style>
