export type SkillDefinition = {
  name: string
  description: string
  unlockAt: number
}

export const languageSkills: Record<string, SkillDefinition[]> = {
  TypeScript: [
    { name: 'Typed JavaScript', description: 'Strongly typed application code', unlockAt: 1 },
    { name: 'React UI', description: 'Component-driven interfaces', unlockAt: 2 },
    { name: 'Browser APIs', description: 'Build interactive web experiences', unlockAt: 3 },
    { name: 'Game Loops', description: 'Real-time browser game logic', unlockAt: 4 },
    { name: 'Production Frontend', description: 'Ship a polished TypeScript app', unlockAt: 5 },
  ],
  Go: [
    { name: 'Go Foundations', description: 'Idiomatic Go application code', unlockAt: 1 },
    { name: 'REST API', description: 'Build HTTP services', unlockAt: 2 },
    { name: 'PostgreSQL', description: 'Persist application data', unlockAt: 3 },
    { name: 'Docker', description: 'Containerize backend services', unlockAt: 4 },
    { name: 'Production Backend', description: 'Ship a complete Go service', unlockAt: 5 },
  ],
  Rust: [
    { name: 'Ownership', description: 'Work with Rust ownership and borrowing', unlockAt: 1 },
    { name: 'Systems', description: 'Build reliable low-level software', unlockAt: 2 },
    { name: 'CLI', description: 'Create useful command-line tools', unlockAt: 3 },
    { name: 'Async Rust', description: 'Build concurrent Rust applications', unlockAt: 4 },
    { name: 'Production Rust', description: 'Ship a robust Rust project', unlockAt: 5 },
  ],
  Python: [
    { name: 'Scripting', description: 'Automate repetitive work', unlockAt: 1 },
    { name: 'Automation', description: 'Turn workflows into tools', unlockAt: 2 },
    { name: 'Data', description: 'Transform and analyze data', unlockAt: 3 },
    { name: 'APIs', description: 'Build useful Python services', unlockAt: 4 },
    { name: 'Production Python', description: 'Ship a complete Python project', unlockAt: 5 },
  ],
}

export function skillsForLanguage(language: string, level: number): SkillDefinition[] {
  const skills = languageSkills[language] ?? [
    { name: 'Foundations', description: `Build with ${language}`, unlockAt: 1 },
    { name: 'Projects', description: `Create projects in ${language}`, unlockAt: 2 },
    { name: 'Tooling', description: `Use the ${language} ecosystem`, unlockAt: 3 },
    { name: 'Production', description: `Ship a ${language} project`, unlockAt: 4 },
    { name: 'Mastery', description: `Go deeper with ${language}`, unlockAt: 5 },
  ]
  return skills.filter((skill) => skill.unlockAt <= level)
}
