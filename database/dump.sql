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
-- Table structure for table `sala`
--

DROP TABLE IF EXISTS `sala`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sala` (
  `id` int NOT NULL AUTO_INCREMENT,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `nome` varchar(255) DEFAULT NULL,
  `ativo` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `nome` (`nome`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sala`
--

LOCK TABLES `sala` WRITE;
/*!40000 ALTER TABLE `sala` DISABLE KEYS */;
INSERT INTO `sala` VALUES (16,'2024-08-03 23:56:49','2024-08-04 20:25:43','oi',0),(17,'2024-08-04 00:56:20','2024-08-04 22:08:06','Amor',0),(18,'2024-08-04 02:50:24','2024-08-04 20:25:45','Ola',0),(19,'2024-08-04 02:50:34','2024-08-04 18:13:06','oi2',0),(20,'2024-08-04 20:09:07','2024-08-04 22:08:06','CArlos',0),(21,'2024-08-04 20:23:23','2024-08-04 20:25:44','a',0),(22,'2024-08-04 20:23:24','2024-08-04 20:25:48','s',0),(23,'2024-08-04 20:23:25','2024-08-04 20:25:49','d',0),(24,'2024-08-04 20:23:27','2024-08-04 20:25:49','f',0),(25,'2024-08-04 20:23:34','2024-08-04 20:25:50','fd',0),(26,'2024-08-04 20:23:36','2024-08-04 20:25:50','dfdf',0),(27,'2024-08-04 20:24:24','2024-08-04 22:08:14','Furacao',0),(28,'2024-08-04 20:25:54','2024-08-04 22:08:57','Atletico',0),(29,'2024-08-04 22:25:20','2024-08-04 22:25:20','salateste',1);
/*!40000 ALTER TABLE `sala` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=792 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuario`
--

LOCK TABLES `usuario` WRITE;
/*!40000 ALTER TABLE `usuario` DISABLE KEYS */;
INSERT INTO `usuario` VALUES (7,'carlos','kadu','2024-08-03 23:56:17','2024-08-03 23:56:17'),(8,'natasha','nathy','2024-08-04 00:08:49','2024-08-04 00:08:49'),(9,'Carlos ribas','Beto','2024-08-04 20:24:53','2024-08-04 20:24:53'),(10,'walter','walter','2024-08-04 20:54:55','2024-08-04 20:54:55'),(11,'usuario1','usuario1','2024-08-05 00:41:13','2024-08-05 00:41:13'),(12,'usuario2','usuario2','2024-08-05 00:41:19','2024-08-05 00:41:19'),(13,'usuario3','usuario3','2024-08-05 00:41:25','2024-08-05 00:41:25'),(14,'usuario4','usuario4','2024-08-05 00:41:32','2024-08-05 00:41:32'),(15,'usuario5','usuario5','2024-08-05 00:41:40','2024-08-05 00:41:40'),(598,'usuario6','usu6','2024-08-05 00:53:00','2024-08-05 00:53:00'),(599,'usuario7','usu7','2024-08-05 00:53:00','2024-08-05 00:53:00'),(600,'usuario8','usu8','2024-08-05 00:53:00','2024-08-05 00:53:00'),(601,'usuario9','usu9','2024-08-05 00:53:00','2024-08-05 00:53:00'),(602,'usuario10','usu10','2024-08-05 00:53:00','2024-08-05 00:53:00'),(603,'usuario11','usu11','2024-08-05 00:53:00','2024-08-05 00:53:00'),(604,'usuario12','usu12','2024-08-05 00:53:00','2024-08-05 00:53:00'),(605,'usuario13','usu13','2024-08-05 00:53:00','2024-08-05 00:53:00'),(606,'usuario14','usu14','2024-08-05 00:53:00','2024-08-05 00:53:00'),(607,'usuario15','usu15','2024-08-05 00:53:00','2024-08-05 00:53:00'),(608,'usuario16','usu16','2024-08-05 00:53:00','2024-08-05 00:53:00'),(609,'usuario17','usu17','2024-08-05 00:53:00','2024-08-05 00:53:00'),(610,'usuario18','usu18','2024-08-05 00:53:00','2024-08-05 00:53:00'),(611,'usuario19','usu19','2024-08-05 00:53:00','2024-08-05 00:53:00'),(612,'usuario20','usu20','2024-08-05 00:53:00','2024-08-05 00:53:00'),(613,'usuario21','usu21','2024-08-05 00:53:00','2024-08-05 00:53:00'),(614,'usuario22','usu22','2024-08-05 00:53:00','2024-08-05 00:53:00'),(615,'usuario23','usu23','2024-08-05 00:53:00','2024-08-05 00:53:00'),(616,'usuario24','usu24','2024-08-05 00:53:00','2024-08-05 00:53:00'),(617,'usuario25','usu25','2024-08-05 00:53:00','2024-08-05 00:53:00'),(618,'usuario26','usu26','2024-08-05 00:53:00','2024-08-05 00:53:00'),(619,'usuario27','usu27','2024-08-05 00:53:00','2024-08-05 00:53:00'),(620,'usuario28','usu28','2024-08-05 00:53:00','2024-08-05 00:53:00'),(621,'usuario29','usu29','2024-08-05 00:53:00','2024-08-05 00:53:00'),(622,'usuario30','usu30','2024-08-05 00:53:00','2024-08-05 00:53:00'),(623,'usuario31','usu31','2024-08-05 00:53:00','2024-08-05 00:53:00'),(624,'usuario32','usu32','2024-08-05 00:53:00','2024-08-05 00:53:00'),(625,'usuario33','usu33','2024-08-05 00:53:00','2024-08-05 00:53:00'),(626,'usuario34','usu34','2024-08-05 00:53:00','2024-08-05 00:53:00'),(627,'usuario35','usu35','2024-08-05 00:53:00','2024-08-05 00:53:00'),(628,'usuario36','usu36','2024-08-05 00:53:00','2024-08-05 00:53:00'),(629,'usuario37','usu37','2024-08-05 00:53:00','2024-08-05 00:53:00'),(630,'usuario38','usu38','2024-08-05 00:53:00','2024-08-05 00:53:00'),(631,'usuario39','usu39','2024-08-05 00:53:00','2024-08-05 00:53:00'),(632,'usuario40','usu40','2024-08-05 00:53:00','2024-08-05 00:53:00'),(633,'usuario41','usu41','2024-08-05 00:53:00','2024-08-05 00:53:00'),(634,'usuario42','usu42','2024-08-05 00:53:00','2024-08-05 00:53:00'),(635,'usuario43','usu43','2024-08-05 00:53:00','2024-08-05 00:53:00'),(636,'usuario44','usu44','2024-08-05 00:53:00','2024-08-05 00:53:00'),(637,'usuario45','usu45','2024-08-05 00:53:00','2024-08-05 00:53:00'),(638,'usuario46','usu46','2024-08-05 00:53:00','2024-08-05 00:53:00'),(639,'usuario47','usu47','2024-08-05 00:53:00','2024-08-05 00:53:00'),(640,'usuario48','usu48','2024-08-05 00:53:00','2024-08-05 00:53:00'),(641,'usuario49','usu49','2024-08-05 00:53:00','2024-08-05 00:53:00'),(642,'usuario50','usu50','2024-08-05 00:53:00','2024-08-05 00:53:00'),(643,'usuario51','usu51','2024-08-05 00:53:00','2024-08-05 00:53:00'),(644,'usuario52','usu52','2024-08-05 00:53:00','2024-08-05 00:53:00'),(645,'usuario53','usu53','2024-08-05 00:53:00','2024-08-05 00:53:00'),(646,'usuario54','usu54','2024-08-05 00:53:00','2024-08-05 00:53:00'),(647,'usuario55','usu55','2024-08-05 00:53:00','2024-08-05 00:53:00'),(648,'usuario56','usu56','2024-08-05 00:53:00','2024-08-05 00:53:00'),(649,'usuario57','usu57','2024-08-05 00:53:00','2024-08-05 00:53:00'),(650,'usuario58','usu58','2024-08-05 00:53:00','2024-08-05 00:53:00'),(651,'usuario59','usu59','2024-08-05 00:53:00','2024-08-05 00:53:00'),(652,'usuario60','usu60','2024-08-05 00:53:00','2024-08-05 00:53:00'),(653,'usuario61','usu61','2024-08-05 00:53:00','2024-08-05 00:53:00'),(654,'usuario62','usu62','2024-08-05 00:53:00','2024-08-05 00:53:00'),(655,'usuario63','usu63','2024-08-05 00:53:00','2024-08-05 00:53:00'),(656,'usuario64','usu64','2024-08-05 00:53:00','2024-08-05 00:53:00'),(657,'usuario65','usu65','2024-08-05 00:53:00','2024-08-05 00:53:00'),(658,'usuario66','usu66','2024-08-05 00:53:00','2024-08-05 00:53:00'),(659,'usuario67','usu67','2024-08-05 00:53:00','2024-08-05 00:53:00'),(660,'usuario68','usu68','2024-08-05 00:53:00','2024-08-05 00:53:00'),(661,'usuario69','usu69','2024-08-05 00:53:00','2024-08-05 00:53:00'),(662,'usuario70','usu70','2024-08-05 00:53:00','2024-08-05 00:53:00'),(663,'usuario71','usu71','2024-08-05 00:53:00','2024-08-05 00:53:00'),(664,'usuario72','usu72','2024-08-05 00:53:00','2024-08-05 00:53:00'),(665,'usuario73','usu73','2024-08-05 00:53:00','2024-08-05 00:53:00'),(666,'usuario74','usu74','2024-08-05 00:53:00','2024-08-05 00:53:00'),(667,'usuario75','usu75','2024-08-05 00:53:00','2024-08-05 00:53:00'),(668,'usuario76','usu76','2024-08-05 00:53:00','2024-08-05 00:53:00'),(669,'usuario77','usu77','2024-08-05 00:53:00','2024-08-05 00:53:00'),(670,'usuario78','usu78','2024-08-05 00:53:00','2024-08-05 00:53:00'),(671,'usuario79','usu79','2024-08-05 00:53:00','2024-08-05 00:53:00'),(672,'usuario80','usu80','2024-08-05 00:53:00','2024-08-05 00:53:00'),(673,'usuario81','usu81','2024-08-05 00:53:00','2024-08-05 00:53:00'),(674,'usuario82','usu82','2024-08-05 00:53:00','2024-08-05 00:53:00'),(675,'usuario83','usu83','2024-08-05 00:53:00','2024-08-05 00:53:00'),(676,'usuario84','usu84','2024-08-05 00:53:00','2024-08-05 00:53:00'),(677,'usuario85','usu85','2024-08-05 00:53:00','2024-08-05 00:53:00'),(678,'usuario86','usu86','2024-08-05 00:53:00','2024-08-05 00:53:00'),(679,'usuario87','usu87','2024-08-05 00:53:00','2024-08-05 00:53:00'),(680,'usuario88','usu88','2024-08-05 00:53:00','2024-08-05 00:53:00'),(681,'usuario89','usu89','2024-08-05 00:53:00','2024-08-05 00:53:00'),(682,'usuario90','usu90','2024-08-05 00:53:00','2024-08-05 00:53:00'),(683,'usuario91','usu91','2024-08-05 00:53:00','2024-08-05 00:53:00'),(684,'usuario92','usu92','2024-08-05 00:53:00','2024-08-05 00:53:00'),(685,'usuario93','usu93','2024-08-05 00:53:00','2024-08-05 00:53:00'),(686,'usuario94','usu94','2024-08-05 00:53:00','2024-08-05 00:53:00'),(687,'usuario95','usu95','2024-08-05 00:53:00','2024-08-05 00:53:00'),(688,'usuario96','usu96','2024-08-05 00:53:00','2024-08-05 00:53:00'),(689,'usuario97','usu97','2024-08-05 00:53:00','2024-08-05 00:53:00'),(690,'usuario98','usu98','2024-08-05 00:53:00','2024-08-05 00:53:00'),(691,'usuario99','usu99','2024-08-05 00:53:00','2024-08-05 00:53:00'),(692,'usuario100','usu100','2024-08-05 00:53:00','2024-08-05 00:53:00'),(693,'usuario101','usu101','2024-08-05 00:53:00','2024-08-05 00:53:00'),(694,'usuario102','usu102','2024-08-05 00:53:00','2024-08-05 00:53:00'),(695,'usuario103','usu103','2024-08-05 00:53:00','2024-08-05 00:53:00'),(696,'usuario104','usu104','2024-08-05 00:53:00','2024-08-05 00:53:00'),(697,'usuario105','usu105','2024-08-05 00:53:00','2024-08-05 00:53:00'),(698,'usuario106','usu106','2024-08-05 00:53:00','2024-08-05 00:53:00'),(699,'usuario107','usu107','2024-08-05 00:53:00','2024-08-05 00:53:00'),(700,'usuario108','usu108','2024-08-05 00:53:00','2024-08-05 00:53:00'),(701,'usuario109','usu109','2024-08-05 00:53:00','2024-08-05 00:53:00'),(702,'usuario110','usu110','2024-08-05 00:53:00','2024-08-05 00:53:00'),(703,'usuario111','usu111','2024-08-05 00:53:00','2024-08-05 00:53:00'),(704,'usuario112','usu112','2024-08-05 00:53:00','2024-08-05 00:53:00'),(705,'usuario113','usu113','2024-08-05 00:53:00','2024-08-05 00:53:00'),(706,'usuario114','usu114','2024-08-05 00:53:00','2024-08-05 00:53:00'),(707,'usuario115','usu115','2024-08-05 00:53:00','2024-08-05 00:53:00'),(708,'usuario116','usu116','2024-08-05 00:53:00','2024-08-05 00:53:00'),(709,'usuario117','usu117','2024-08-05 00:53:00','2024-08-05 00:53:00'),(710,'usuario118','usu118','2024-08-05 00:53:00','2024-08-05 00:53:00'),(711,'usuario119','usu119','2024-08-05 00:53:00','2024-08-05 00:53:00'),(712,'usuario120','usu120','2024-08-05 00:53:00','2024-08-05 00:53:00'),(713,'usuario121','usu121','2024-08-05 00:53:00','2024-08-05 00:53:00'),(714,'usuario122','usu122','2024-08-05 00:53:00','2024-08-05 00:53:00'),(715,'usuario123','usu123','2024-08-05 00:53:00','2024-08-05 00:53:00'),(716,'usuario124','usu124','2024-08-05 00:53:00','2024-08-05 00:53:00'),(717,'usuario125','usu125','2024-08-05 00:53:00','2024-08-05 00:53:00'),(718,'usuario126','usu126','2024-08-05 00:53:00','2024-08-05 00:53:00'),(719,'usuario127','usu127','2024-08-05 00:53:00','2024-08-05 00:53:00'),(720,'usuario128','usu128','2024-08-05 00:53:00','2024-08-05 00:53:00'),(721,'usuario129','usu129','2024-08-05 00:53:00','2024-08-05 00:53:00'),(722,'usuario130','usu130','2024-08-05 00:53:00','2024-08-05 00:53:00'),(723,'usuario131','usu131','2024-08-05 00:53:00','2024-08-05 00:53:00'),(724,'usuario132','usu132','2024-08-05 00:53:00','2024-08-05 00:53:00'),(725,'usuario133','usu133','2024-08-05 00:53:00','2024-08-05 00:53:00'),(726,'usuario134','usu134','2024-08-05 00:53:00','2024-08-05 00:53:00'),(727,'usuario135','usu135','2024-08-05 00:53:00','2024-08-05 00:53:00'),(728,'usuario136','usu136','2024-08-05 00:53:00','2024-08-05 00:53:00'),(729,'usuario137','usu137','2024-08-05 00:53:00','2024-08-05 00:53:00'),(730,'usuario138','usu138','2024-08-05 00:53:00','2024-08-05 00:53:00'),(731,'usuario139','usu139','2024-08-05 00:53:00','2024-08-05 00:53:00'),(732,'usuario140','usu140','2024-08-05 00:53:00','2024-08-05 00:53:00'),(733,'usuario141','usu141','2024-08-05 00:53:00','2024-08-05 00:53:00'),(734,'usuario142','usu142','2024-08-05 00:53:00','2024-08-05 00:53:00'),(735,'usuario143','usu143','2024-08-05 00:53:00','2024-08-05 00:53:00'),(736,'usuario144','usu144','2024-08-05 00:53:00','2024-08-05 00:53:00'),(737,'usuario145','usu145','2024-08-05 00:53:00','2024-08-05 00:53:00'),(738,'usuario146','usu146','2024-08-05 00:53:00','2024-08-05 00:53:00'),(739,'usuario147','usu147','2024-08-05 00:53:00','2024-08-05 00:53:00'),(740,'usuario148','usu148','2024-08-05 00:53:00','2024-08-05 00:53:00'),(741,'usuario149','usu149','2024-08-05 00:53:00','2024-08-05 00:53:00'),(742,'usuario150','usu150','2024-08-05 00:53:00','2024-08-05 00:53:00'),(743,'usuario151','usu151','2024-08-05 00:53:00','2024-08-05 00:53:00'),(744,'usuario152','usu152','2024-08-05 00:53:00','2024-08-05 00:53:00'),(745,'usuario153','usu153','2024-08-05 00:53:00','2024-08-05 00:53:00'),(746,'usuario154','usu154','2024-08-05 00:53:00','2024-08-05 00:53:00'),(747,'usuario155','usu155','2024-08-05 00:53:00','2024-08-05 00:53:00'),(748,'usuario156','usu156','2024-08-05 00:53:00','2024-08-05 00:53:00'),(749,'usuario157','usu157','2024-08-05 00:53:00','2024-08-05 00:53:00'),(750,'usuario158','usu158','2024-08-05 00:53:00','2024-08-05 00:53:00'),(751,'usuario159','usu159','2024-08-05 00:53:00','2024-08-05 00:53:00'),(752,'usuario160','usu160','2024-08-05 00:53:00','2024-08-05 00:53:00'),(753,'usuario161','usu161','2024-08-05 00:53:00','2024-08-05 00:53:00'),(754,'usuario162','usu162','2024-08-05 00:53:00','2024-08-05 00:53:00'),(755,'usuario163','usu163','2024-08-05 00:53:00','2024-08-05 00:53:00'),(756,'usuario164','usu164','2024-08-05 00:53:00','2024-08-05 00:53:00'),(757,'usuario165','usu165','2024-08-05 00:53:00','2024-08-05 00:53:00'),(758,'usuario166','usu166','2024-08-05 00:53:00','2024-08-05 00:53:00'),(759,'usuario167','usu167','2024-08-05 00:53:00','2024-08-05 00:53:00'),(760,'usuario168','usu168','2024-08-05 00:53:00','2024-08-05 00:53:00'),(761,'usuario169','usu169','2024-08-05 00:53:00','2024-08-05 00:53:00'),(762,'usuario170','usu170','2024-08-05 00:53:00','2024-08-05 00:53:00'),(763,'usuario171','usu171','2024-08-05 00:53:00','2024-08-05 00:53:00'),(764,'usuario172','usu172','2024-08-05 00:53:00','2024-08-05 00:53:00'),(765,'usuario173','usu173','2024-08-05 00:53:00','2024-08-05 00:53:00'),(766,'usuario174','usu174','2024-08-05 00:53:00','2024-08-05 00:53:00'),(767,'usuario175','usu175','2024-08-05 00:53:00','2024-08-05 00:53:00'),(768,'usuario176','usu176','2024-08-05 00:53:00','2024-08-05 00:53:00'),(769,'usuario177','usu177','2024-08-05 00:53:00','2024-08-05 00:53:00'),(770,'usuario178','usu178','2024-08-05 00:53:00','2024-08-05 00:53:00'),(771,'usuario179','usu179','2024-08-05 00:53:00','2024-08-05 00:53:00'),(772,'usuario180','usu180','2024-08-05 00:53:00','2024-08-05 00:53:00'),(773,'usuario181','usu181','2024-08-05 00:53:00','2024-08-05 00:53:00'),(774,'usuario182','usu182','2024-08-05 00:53:00','2024-08-05 00:53:00'),(775,'usuario183','usu183','2024-08-05 00:53:00','2024-08-05 00:53:00'),(776,'usuario184','usu184','2024-08-05 00:53:00','2024-08-05 00:53:00'),(777,'usuario185','usu185','2024-08-05 00:53:00','2024-08-05 00:53:00'),(778,'usuario186','usu186','2024-08-05 00:53:00','2024-08-05 00:53:00'),(779,'usuario187','usu187','2024-08-05 00:53:00','2024-08-05 00:53:00'),(780,'usuario188','usu188','2024-08-05 00:53:00','2024-08-05 00:53:00'),(781,'usuario189','usu189','2024-08-05 00:53:00','2024-08-05 00:53:00'),(782,'usuario190','usu190','2024-08-05 00:53:00','2024-08-05 00:53:00'),(783,'usuario191','usu191','2024-08-05 00:53:00','2024-08-05 00:53:00'),(784,'usuario192','usu192','2024-08-05 00:53:00','2024-08-05 00:53:00'),(785,'usuario193','usu193','2024-08-05 00:53:00','2024-08-05 00:53:00'),(786,'usuario194','usu194','2024-08-05 00:53:00','2024-08-05 00:53:00'),(787,'usuario195','usu195','2024-08-05 00:53:00','2024-08-05 00:53:00'),(788,'usuario196','usu196','2024-08-05 00:53:00','2024-08-05 00:53:00'),(789,'usuario197','usu197','2024-08-05 00:53:00','2024-08-05 00:53:00'),(790,'usuario198','usu198','2024-08-05 00:53:00','2024-08-05 00:53:00'),(791,'usuario199','usu199','2024-08-05 00:53:00','2024-08-05 00:53:00');
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

--
-- Table structure for table `usuariosala`
--

DROP TABLE IF EXISTS `usuariosala`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usuariosala` (
  `id` int NOT NULL AUTO_INCREMENT,
  `idusuario` int NOT NULL,
  `idsala` int NOT NULL,
  `datacriacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `dataatualizacao` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `online` tinyint(1) DEFAULT '0',
  `ativo` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idusuario_2` (`idusuario`,`idsala`),
  KEY `idsala` (`idsala`),
  KEY `idusuario` (`idusuario`),
  CONSTRAINT `usuariosala_ibfk_1` FOREIGN KEY (`idsala`) REFERENCES `sala` (`id`),
  CONSTRAINT `usuariosala_ibfk_2` FOREIGN KEY (`idusuario`) REFERENCES `usuario` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usuariosala`
--

LOCK TABLES `usuariosala` WRITE;
/*!40000 ALTER TABLE `usuariosala` DISABLE KEYS */;
INSERT INTO `usuariosala` VALUES (9,7,16,'2024-08-03 23:56:49','2024-08-04 20:25:43',0,0),(13,8,16,'2024-08-04 00:08:54','2024-08-04 20:24:20',0,0),(14,7,17,'2024-08-04 00:56:20','2024-08-04 22:08:06',0,0),(15,8,17,'2024-08-04 00:56:36','2024-08-04 20:24:20',0,0),(21,7,18,'2024-08-04 02:50:24','2024-08-04 20:25:45',0,0),(22,7,19,'2024-08-04 02:50:34','2024-08-04 18:13:06',0,0),(25,7,20,'2024-08-04 20:09:07','2024-08-04 22:08:06',0,0),(26,7,21,'2024-08-04 20:23:23','2024-08-04 20:25:44',0,0),(27,7,22,'2024-08-04 20:23:24','2024-08-04 20:25:48',0,0),(28,7,23,'2024-08-04 20:23:25','2024-08-04 20:25:49',0,0),(29,7,24,'2024-08-04 20:23:27','2024-08-04 20:25:49',0,0),(30,7,25,'2024-08-04 20:23:34','2024-08-04 20:25:50',0,0),(31,7,26,'2024-08-04 20:23:36','2024-08-04 20:25:50',0,0),(32,8,27,'2024-08-04 20:24:24','2024-08-04 22:08:14',0,0),(33,7,28,'2024-08-04 20:25:54','2024-08-04 22:08:05',0,0),(34,9,28,'2024-08-04 20:53:10','2024-08-04 22:08:21',0,0),(35,10,28,'2024-08-04 20:55:04','2024-08-04 22:08:57',0,0),(36,7,29,'2024-08-04 22:25:20','2024-08-04 22:25:20',0,1),(37,8,29,'2024-08-04 22:27:07','2024-08-04 22:27:07',0,1),(38,9,29,'2024-08-04 22:27:07','2024-08-04 22:27:07',0,1),(39,10,29,'2024-08-04 22:27:07','2024-08-04 22:27:07',0,1);
/*!40000 ALTER TABLE `usuariosala` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-05  0:58:40
