export type CollectionSkill = {
  language: string
  skills: string[]
}

export const collectionSkills: CollectionSkill[] = [
  { language: 'TypeScript', skills: ['React UI', 'Browser APIs', 'Game loops'] },
  { language: 'JavaScript', skills: ['DOM', 'Game logic', 'Web APIs'] },
  { language: 'Go', skills: ['REST API', 'PostgreSQL', 'Docker'] },
  { language: 'Rust', skills: ['Systems', 'Ownership', 'CLI'] },
  { language: 'Python', skills: ['Automation', 'Data', 'Scripting'] },
]

export function skillsForLanguage(language: string): string[] {
  return collectionSkills.find((item) => item.language === language)?.skills ?? ['First project', 'Keep building', 'Keep learning']
}
