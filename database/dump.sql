-- MySQL dump 10.13  Distrib 9.0.0, for Linux (x86_64)
--
-- Host: localhost    Database: chat
-- ------------------------------------------------------
-- Server version	9.0.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `canal`
--

DROP TABLE IF EXISTS `canal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `canal` (
  `id` int NOT NULL AUTO_INCREMENT,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `nome` varchar(255) DEFAULT NULL,
  `online` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `nome` (`nome`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `canal`
--

LOCK TABLES `canal` WRITE;
/*!40000 ALTER TABLE `canal` DISABLE KEYS */;
INSERT INTO `canal` VALUES (1,'2024-06-30 14:58:16',NULL,'teste',0),(6,'2024-07-20 13:39:43','2024-07-20 13:39:43','canal da nathy e do kady',0),(7,'2024-07-20 13:43:07','2024-07-20 22:55:33','crypto',0),(11,'2024-07-20 21:32:31','2024-07-20 21:32:31','andromeda',0),(12,'2024-07-20 21:33:43','2024-07-20 21:33:43','gargantua',0);
/*!40000 ALTER TABLE `canal` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `canalconexoes`
--

DROP TABLE IF EXISTS `canalconexoes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `canalconexoes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `idcanal` int NOT NULL,
  `idconexao` char(36) NOT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idcanal` (`idcanal`),
  KEY `idconexao` (`idconexao`),
  CONSTRAINT `canalconexoes_ibfk_1` FOREIGN KEY (`idcanal`) REFERENCES `canal` (`id`),
  CONSTRAINT `canalconexoes_ibfk_2` FOREIGN KEY (`idconexao`) REFERENCES `conexao` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `canalconexoes`
--

LOCK TABLES `canalconexoes` WRITE;
/*!40000 ALTER TABLE `canalconexoes` DISABLE KEYS */;
/*!40000 ALTER TABLE `canalconexoes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `canaliteracao`
--

DROP TABLE IF EXISTS `canaliteracao`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `canaliteracao` (
  `id` int NOT NULL AUTO_INCREMENT,
  `mensagem` text,
  `anexo` varchar(255) DEFAULT NULL,
  `idcanalconexao` int NOT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `canaliteracao`
--

LOCK TABLES `canaliteracao` WRITE;
/*!40000 ALTER TABLE `canaliteracao` DISABLE KEYS */;
/*!40000 ALTER TABLE `canaliteracao` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conexao`
--

DROP TABLE IF EXISTS `conexao`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `conexao` (
  `id` char(36) NOT NULL DEFAULT (uuid()),
  `nome` varchar(100) NOT NULL,
  `idsocket` varchar(36) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conexao`
--

LOCK TABLES `conexao` WRITE;
/*!40000 ALTER TABLE `conexao` DISABLE KEYS */;
/*!40000 ALTER TABLE `conexao` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usuario`
--

DROP TABLE IF EXISTS `usuario`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuario` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nome` varchar(255) DEFAULT NULL,
  `apelido` varchar(8) DEFAULT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `apelido` (`apelido`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuario`
--

LOCK TABLES `usuario` WRITE;
/*!40000 ALTER TABLE `usuario` DISABLE KEYS */;
INSERT INTO `usuario` VALUES (1,'Kaduhodi user name','Kaduhodi','2024-07-08 12:29:29','2024-07-08 12:29:29'),(2,'carlosAl user name','carlosAl','2024-07-08 12:32:03','2024-07-08 12:32:03'),(3,'Arnilloy user name','Arnilloy','2024-07-08 16:57:11','2024-07-08 16:57:11');
/*!40000 ALTER TABLE `usuario` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usuariocanal`
--

DROP TABLE IF EXISTS `usuariocanal`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuariocanal` (
  `id` int NOT NULL AUTO_INCREMENT,
  `idusuario` int NOT NULL,
  `idcanal` int NOT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `online` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idusuario_2` (`idusuario`,`idcanal`),
  KEY `idcanal` (`idcanal`),
  KEY `idusuario` (`idusuario`),
  CONSTRAINT `usuariocanal_ibfk_1` FOREIGN KEY (`idcanal`) REFERENCES `canal` (`id`),
  CONSTRAINT `usuariocanal_ibfk_2` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuariocanal`
--

LOCK TABLES `usuariocanal` WRITE;
/*!40000 ALTER TABLE `usuariocanal` DISABLE KEYS */;
INSERT INTO `usuariocanal` VALUES (1,1,7,'2024-07-21 00:17:13','2024-07-21 01:44:46',0);
/*!40000 ALTER TABLE `usuariocanal` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usuarioconexao`
--

DROP TABLE IF EXISTS `usuarioconexao`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuarioconexao` (
  `id` int NOT NULL AUTO_INCREMENT,
  `idusuario` int NOT NULL,
  `idconexao` char(36) NOT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idusuario` (`idusuario`),
  KEY `idconexao` (`idconexao`),
  CONSTRAINT `usuarioconexao_ibfk_1` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`),
  CONSTRAINT `usuarioconexao_ibfk_2` FOREIGN KEY (`idconexao`) REFERENCES `conexao` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuarioconexao`
--

LOCK TABLES `usuarioconexao` WRITE;
/*!40000 ALTER TABLE `usuarioconexao` DISABLE KEYS */;
/*!40000 ALTER TABLE `usuarioconexao` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-07-21  1:53:57

