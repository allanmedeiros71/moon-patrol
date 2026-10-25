# Projeto: Moon Patrol — Releitura Moderna (Linux)

## Objetivo

Recriar o jogo arcade Moon Patrol (Irem, 1982) como releitura moderna,
preservando fielmente a dinâmica e a mecânica originais, com visual e som novos.
Não usar código, gráficos ou sons do original (todos os assets são próprios).

## Stack

- Linguagem: Go (versão estável mais recente)
- Engine: Ebitengine
- Alvo: Linux (x86_64); build reproduzível via Makefile
- Repositório novo, criado do zero

## Requisitos de tela

- Tela cheia por padrão, alternável para janela (F11 / Alt+Enter)
- Redimensionamento para qualquer resolução, mantendo proporção 16:9
  (letterbox quando necessário), lógica de jogo independente da resolução
- 60 FPS com lógica em passo fixo

## Mecânica (fiel ao original)

- Buggy lunar anda da esquerda para a direita; jogador controla velocidade
  (acelerar/frear) e pulo
- Tiro duplo: para frente (destrói pedras/obstáculos) e para cima (inimigos aéreos)
- Obstáculos: crateras, pedras, minas, crateras criadas por bombas
- Inimigos: OVNIs/naves que atiram e bombardeiam, tanques, aviões rasantes
- Rolagem horizontal com parallax de múltiplas camadas (montanhas, cidade, solo)
- Percurso com checkpoints A–Z; bônus por tempo ao chegar no checkpoint
- Dois percursos (iniciante e campeão), vidas, continuação a partir do último checkpoint
- Pontuação conforme o original (inimigos, pedras, pulos sobre obstáculos)

## Controles

- Teclado (configurável) e gamepad
- Padrão: ←/→ velocidade, ↑ ou Z pulo, X/Espaço tiro

## Visual e áudio modernos

- Arte nova em alta resolução (sprites ou vetorial), estilo retrô-futurista
- Partículas (explosões, poeira), iluminação/brilho, transições suaves
- Trilha sonora nova inspirada no tema original (sem copiar) e efeitos sonoros

## Interface

- Menu inicial, pausa, opções (vídeo, áudio, controles), tela de game over
- Recordes locais persistidos em ~/.local/share/moonpatrol/
- HUD: pontuação, recorde, vidas, checkpoint atual, tempo

## Arquitetura e qualidade

- Pacotes separados: game (estados), entities, physics, level, render, audio, input
- Fases definidas em arquivos de dados (JSON/YAML), não hardcoded
- Testes unitários para física, colisão e pontuação
- README com instruções de build/execução

## Entregas por etapa

1. Base: projeto, loop, tela cheia redimensionável, input
2. Mecânica: buggy, pulo, parallax, crateras/pedras, tiros
3. Inimigos e fases: OVNIs, bombas, tanques, checkpoints, placar, vidas
4. Acabamento: arte/som modernos, menus, recordes, pacote Linux (AppImage ou .tar.gz)

Ao fim de cada etapa: o jogo deve compilar e rodar; aguardar revisão antes da próxima.
