-- Script para criar o banco de dados MySQL
-- Execute este script no seu MySQL para preparar o ambiente

CREATE DATABASE IF NOT EXISTS app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE app_db;

-- Verificar se as tabelas foram criadas (serão criadas automaticamente pelo GORM)
SHOW TABLES;