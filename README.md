# 🐉 AI-Powered RPG Engine

A modular Go RPG engine where the engine controls the gameplay and AI generates dynamic storytelling.

This project combines a deterministic game engine with an LLM (Ollama) to generate dynamic adventures while keeping all gameplay rules under engine control.

Instead of allowing the AI to control the game logic, the engine validates and manages every important mechanic such as combat, progression, monsters, inventory and story flow.

---

## ✨ Features

- AI-generated adventure scenes using Ollama
- Turn-based combat system
- Class progression
- Equipment system
- Inventory management
- Status effects (Debuffs)
- PostgreSQL save system
- JSONB player persistence
- Story progression
- Biome-based monster selection
- Deterministic game rules
- Telegram-ready architecture

---

# Adventure Structure

Each adventure contains **15 scenes**.

```
Scene 1
│
├── Scene 2
├── Scene 3 (Required Combat)
├── Scene 4
├── Scene 5
│
├── Scene 6 (Required Combat)
├── Scene 7
├── Scene 8
│
├── Scene 9 (Required Combat)
├── Scene 10
│
├── Scene 11
├── Scene 12 (Required Combat)
├── Scene 13
├── Scene 14
│
└── Scene 15 (Final Boss)
```

---

# Regions

## 🌲 Dark Forest

Scenes **1-5**

Possible monsters:

- Goblin
- Orc

---

## 🏛 Ancient Ruins

Scenes **6-10**

Possible monsters:

- Orc
- Guardian

---

## ⛪ Sanctuary

Scenes **11-14**

Possible monsters:

- Dragon Statue

Scene **15**

- Dragon (Final Boss)

---

# Combat Rules

| Scene | Combat |
|-------|---------|
| 1 | Never |
| Every 3rd scene | Required |
| Remaining scenes | 50% chance |

The engine decides if combat happens.

The AI only creates the narrative.

---

# AI Responsibilities

The LLM is responsible for:

- Writing scenes
- Writing choices
- Selecting one allowed monster when combat is enabled

The AI is **not allowed** to:

- Create monsters
- Give rewards
- Modify player stats
- Change progression
- Invent equipment
- Skip engine rules

---

# Engine Responsibilities

The engine controls:

- Story progression
- Combat
- Inventory
- Equipment
- Skills
- XP
- Gold
- Drops
- Regions
- Save/Load
- Validation
- Boss spawning

---

# Technologies

- Go
- PostgreSQL
- Ollama
- Telegram Bot API 

---

# Project Structure

```
internal/

├── app/
├── config/
├── database/
├── domain/
├── engine/
├── handlers/
├── ollama/
├── repository/
├── server/
└── services/
```

# Philosophy

The goal of this project is **not** to let AI run the game.

Instead:

- The engine controls the rules.
- The AI creates the storytelling.

This approach guarantees consistency, balance and deterministic gameplay while still allowing every adventure to feel unique.

---

## License

MIT
