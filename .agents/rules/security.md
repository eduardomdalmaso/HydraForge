# 🔒 Diretrizes Invioláveis de Segurança & Hardening (HydraForge)

Este documento estabelece as diretrizes invioláveis de segurança, autenticação de jobs de treinamento YOLO, proteção de modelos exportados TensorRT/ONNX e isolamento multi-tenant para o HydraForge.

---

## 1. Isolamento Multi-Tenant em Datasets & Treinamento
- **Isolamento de Datasets:** O worker de treinamento na RTX 5090 consome apenas datasets versionados pertencentes ao `tenant_id` solicitante.
- **Isolamento de Checkpoints e Pesos (.engine / .pt):** Modelos compilados são persistidos em diretórios isolados por tenant no MinIO S3 ou filesystem com controle rigoroso de leitura.

---

## 2. Permissão no Servidor (Backend RBAC)
- **Controle de Execução de Treino:** Apenas usuários autenticados com papel `admin` ou `researcher` podem despachar jobs de treinamento intensivo na GPU.

---

## 3. Prevenção a IDOR & Download de Modelos
- **Validação de Modelo por ID:** Endpoints de download de modelos compilados (`/api/v1/models/{id}/export`) validam a propriedade do modelo pelo tenant do token antes da emissão de links de download.

---

## 4. Chaves & Segredos (Zero Hardcoded Secrets)
- **Tokens do Hugging Face & GitHub Releases:** Chaves de API para publicação de modelos devem ser injetadas exclusivamente via variáveis de ambiente (`HF_TOKEN`, `GITHUB_TOKEN`), nunca commitadas no repositório.

---

## 5. Sanitização de Inputs & Prevenção de Injeção de Comandos
- **Execução de Processos Python Ultralytics:** Parâmetros de treino (`epochs`, `imgsz`, `batch`, `weights`) devem ser estritamente tipados e validados por limites numéricos de domínio para impedir injeção de comandos de shell via CLI.
