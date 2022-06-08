module.exports = {
  'sort-imports': ['warn', { ignoreDeclarationSort: true }],
  'import/first': 'warn',
  'import/order': [
    'warn',
    {
      groups: [
        'builtin',
        'external',
        'internal',
        'parent',
        'sibling',
        'index',
        'object',
        'type',
      ],
      'newlines-between': 'always',
      alphabetize: {
        order: 'asc',
        caseInsensitive: true,
      },
    },
  ],
  'import/newline-after-import': 'warn',
}
