/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'Data Connector',
  tagline: 'Data Connector for Google Sheets',
  url: 'https://dataconnector.app/docs',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.png',
  organizationName: 'brentadamson', // Usually your GitHub org/user name.
  projectName: 'dataconnector', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Data Connector',
      logo: {
        alt: 'Data Connector Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {
          href: 'https://github.com/brentadamson/dataconnector',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: 'docs/',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Twitter',
              href: 'https://twitter.com/data_connector',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/brentadamson/dataconnector',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Data Connector. Docs built with Docusaurus.`,
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:
            'https://github.com/brentadamson/dataconnector/edit/master/website/',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          editUrl:
            'https://github.com/brentadamson/dataconnector/edit/master/website/blog/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
