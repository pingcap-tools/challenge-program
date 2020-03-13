-- MySQL dump 10.13  Distrib 5.7.29, for Linux (x86_64)
--
-- Host: 172.17.0.1    Database: community
-- ------------------------------------------------------
-- Server version	5.7.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(255) DEFAULT NULL,
  `repo` varchar(255) DEFAULT NULL,
  `comment_id` int(11) DEFAULT NULL,
  `comment_type` varchar(128) DEFAULT NULL,
  `pull_number` int(11) DEFAULT NULL,
  `body` text,
  `user` varchar(255) DEFAULT NULL,
  `url` varchar(1023) DEFAULT NULL,
  `association` varchar(255) DEFAULT NULL,
  `relation` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_comments_pull_number` (`pull_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `picks`
--

DROP TABLE IF EXISTS `picks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `picks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `season` int(11) DEFAULT NULL,
  `owner` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `repo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `task_id` int(11) DEFAULT '0',
  `teamID` int(11) DEFAULT NULL,
  `user` varchar(255) DEFAULT NULL,
  `pull_number` int(11) DEFAULT '0',
  `score` int(11) NOT NULL,
  `status` varchar(128) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `closed_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `picks`
--

LOCK TABLES `picks` WRITE;
/*!40000 ALTER TABLE `picks` DISABLE KEYS */;
/*!40000 ALTER TABLE `picks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pulls`
--

DROP TABLE IF EXISTS `pulls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pulls` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `owner` varchar(255) DEFAULT NULL,
  `repo` varchar(255) DEFAULT NULL,
  `pull_number` int(11) DEFAULT NULL,
  `title` text,
  `body` text,
  `user` varchar(255) DEFAULT NULL,
  `association` varchar(255) DEFAULT NULL,
  `relation` varchar(255) DEFAULT NULL,
  `label` text,
  `status` varchar(128) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `closed_at` timestamp NULL DEFAULT NULL,
  `merged_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_pull_number` (`pull_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pulls`
--

LOCK TABLES `pulls` WRITE;
/*!40000 ALTER TABLE `pulls` DISABLE KEYS */;
/*!40000 ALTER TABLE `pulls` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reviewers`
--

DROP TABLE IF EXISTS `reviewers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reviewers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `github_id` varchar(255) DEFAULT NULL,
  `organization` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `slack_id` varchar(255) DEFAULT NULL,
  `group_name` varchar(255) DEFAULT NULL,
  `role` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reviewers`
--

LOCK TABLES `reviewers` WRITE;
/*!40000 ALTER TABLE `reviewers` DISABLE KEYS */;
/*!40000 ALTER TABLE `reviewers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `taskgroups`
--

DROP TABLE IF EXISTS `taskgroups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `taskgroups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `season` int(11) DEFAULT NULL,
  `owner` varchar(255) DEFAULT NULL,
  `repo` varchar(255) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `issue_number` int(11) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `bonus` int(11) DEFAULT '0',
  `vote` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `taskgroups`
--

LOCK TABLES `taskgroups` WRITE;
/*!40000 ALTER TABLE `taskgroups` DISABLE KEYS */;
/*!40000 ALTER TABLE `taskgroups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `season` int(11) DEFAULT NULL,
  `complete_user` varchar(255) DEFAULT NULL,
  `complete_team` int(11) DEFAULT NULL,
  `owner` varchar(255) DEFAULT NULL,
  `repo` varchar(255) DEFAULT NULL,
  `title` varchar(2047) DEFAULT NULL,
  `issue_number` int(11) NOT NULL DEFAULT '0',
  `pull_number` int(11) DEFAULT NULL,
  `level` varchar(255) DEFAULT NULL,
  `min_score` int(11) DEFAULT '0',
  `score` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expired` varchar(255) DEFAULT NULL,
  `taskgroup_id` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teams`
--

DROP TABLE IF EXISTS `teams`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `teams` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `season` int(11) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `issue_url` varchar(1023) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teams`
--

LOCK TABLES `teams` WRITE;
/*!40000 ALTER TABLE `teams` DISABLE KEYS */;
/*!40000 ALTER TABLE `teams` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `season` int(11) DEFAULT NULL,
  `user` varchar(255) DEFAULT NULL,
  `email` varchar(1023) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `team_id` int(11) DEFAULT NULL,
  `issue_url` varchar(1023) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `leader` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `views`
--

DROP TABLE IF EXISTS `views`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `views` (
  `viewname` varchar(100) DEFAULT NULL,
  `displayname` varchar(200) DEFAULT NULL,
  `description` text,
  `visible` tinyint(1) DEFAULT NULL,
  UNIQUE KEY `viewname` (`viewname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `views`
--

LOCK TABLES `views` WRITE;
/*!40000 ALTER TABLE `views` DISABLE KEYS */;
/*!40000 ALTER TABLE `views` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-03-13 18:05:19
