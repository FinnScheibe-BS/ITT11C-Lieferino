
-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Erstellungszeit: 30. Jun 2026 um 09:19
-- Server-Version: 10.4.28-MariaDB
-- PHP-Version: 8.2.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Datenbank: `aup`
--

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Bestellung`
--

CREATE TABLE `Bestellung` (
  `PK_ID_Bestellung` int(11) NOT NULL,
  `FK_ID_Restaurant` int(11) NOT NULL,
  `FK_ID_Kunde` int(11) NOT NULL,
  `FK_ID_Gericht` int(11) NOT NULL,
  `Datum` date NOT NULL,
  `Uhrzeit` time NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Bewertung`
--

CREATE TABLE `Bewertung` (
  `FK_ID_Kunde` int(11) NOT NULL,
  `FK_ID_Restaurant` int(11) NOT NULL,
  `Sterne` int(11) NOT NULL,
  `Text` text DEFAULT NULL,
  `Uhrzeit` time DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Gericht`
--

CREATE TABLE `Gericht` (
  `PK_ID_Gericht` int(11) NOT NULL,
  `FK_ID_Restaurant` int(11) DEFAULT NULL,
  `Name` varchar(20) NOT NULL,
  `Vegetarisch` bit(1) NOT NULL,
  `Vegan` bit(1) NOT NULL,
  `Preis` float NOT NULL,
  `FK_ID_Kategorie` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Gericht_Zutat`
--

CREATE TABLE `Gericht_Zutat` (
  `FK_ID_Gericht` int(11) NOT NULL,
  `FK_ID_Zutat` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Kategorie`
--

CREATE TABLE `Kategorie` (
  `PK_ID_Kategorie` int(11) NOT NULL,
  `Bezeichnung` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Kunde`
--

CREATE TABLE `Kunde` (
  `PK_ID_Kunde` int(11) NOT NULL,
  `Vorname` varchar(50) NOT NULL,
  `Nachname` varchar(50) NOT NULL,
  `Adresse` varchar(100) DEFAULT NULL,
  `Telefonnummer` varchar(20) DEFAULT NULL,
  `Email_Adresse` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Restaurant`
--

CREATE TABLE `Restaurant` (
  `PK_ID_Restaurant` int(11) NOT NULL,
  `Name` varchar(20) NOT NULL,
  `Adresse` varchar(80) NOT NULL,
  `Nationalitaet` varchar(20) DEFAULT NULL,
  `Preisklasse` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Daten für Tabelle `Restaurant`
--

INSERT INTO `Restaurant` (`PK_ID_Restaurant`, `Name`, `Adresse`, `Nationalitaet`, `Preisklasse`) VALUES
(1, 'Test', 'TestStraße', 'De', 1);

-- --------------------------------------------------------

--
-- Tabellenstruktur für Tabelle `Zutat`
--

CREATE TABLE `Zutat` (
  `PK_ID_Zutat` int(11) NOT NULL,
  `Name` varchar(50) NOT NULL,
  `Kalorien` int(11) DEFAULT NULL,
  `Naehrwerte` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indizes der exportierten Tabellen
--

--
-- Indizes für die Tabelle `Bestellung`
--
ALTER TABLE `Bestellung`
  ADD PRIMARY KEY (`PK_ID_Bestellung`),
  ADD KEY `FK_ID_Restaurant` (`FK_ID_Restaurant`),
  ADD KEY `FK_ID_Kunde` (`FK_ID_Kunde`),
  ADD KEY `FK_ID_Gericht` (`FK_ID_Gericht`);

--
-- Indizes für die Tabelle `Bewertung`
--
ALTER TABLE `Bewertung`
  ADD PRIMARY KEY (`FK_ID_Kunde`,`FK_ID_Restaurant`),
  ADD KEY `FK_ID_Restaurant` (`FK_ID_Restaurant`);

--
-- Indizes für die Tabelle `Gericht`
--
ALTER TABLE `Gericht`
  ADD PRIMARY KEY (`PK_ID_Gericht`),
  ADD KEY `FK_ID_Restaurant` (`FK_ID_Restaurant`),
  ADD KEY `FK_Gericht_Kategorie` (`FK_ID_Kategorie`);

--
-- Indizes für die Tabelle `Gericht_Zutat`
--
ALTER TABLE `Gericht_Zutat`
  ADD PRIMARY KEY (`FK_ID_Gericht`,`FK_ID_Zutat`),
  ADD KEY `FK_ID_Zutat` (`FK_ID_Zutat`);

--
-- Indizes für die Tabelle `Kategorie`
--
ALTER TABLE `Kategorie`
  ADD PRIMARY KEY (`PK_ID_Kategorie`);

--
-- Indizes für die Tabelle `Kunde`
--
ALTER TABLE `Kunde`
  ADD PRIMARY KEY (`PK_ID_Kunde`),
  ADD UNIQUE KEY `Email_Adresse` (`Email_Adresse`);

--
-- Indizes für die Tabelle `Restaurant`
--
ALTER TABLE `Restaurant`
  ADD PRIMARY KEY (`PK_ID_Restaurant`);

--
-- Indizes für die Tabelle `Zutat`
--
ALTER TABLE `Zutat`
  ADD PRIMARY KEY (`PK_ID_Zutat`);

--
-- AUTO_INCREMENT für exportierte Tabellen
--

--
-- AUTO_INCREMENT für Tabelle `Bestellung`
--
ALTER TABLE `Bestellung`
  MODIFY `PK_ID_Bestellung` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT für Tabelle `Gericht`
--
ALTER TABLE `Gericht`
  MODIFY `PK_ID_Gericht` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT für Tabelle `Kategorie`
--
ALTER TABLE `Kategorie`
  MODIFY `PK_ID_Kategorie` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT für Tabelle `Kunde`
--
ALTER TABLE `Kunde`
  MODIFY `PK_ID_Kunde` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT für Tabelle `Restaurant`
--
ALTER TABLE `Restaurant`
  MODIFY `PK_ID_Restaurant` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT für Tabelle `Zutat`
--
ALTER TABLE `Zutat`
  MODIFY `PK_ID_Zutat` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints der exportierten Tabellen
--

--
-- Constraints der Tabelle `Bestellung`
--
ALTER TABLE `Bestellung`
  ADD CONSTRAINT `bestellung_ibfk_1` FOREIGN KEY (`FK_ID_Restaurant`) REFERENCES `Restaurant` (`PK_ID_Restaurant`),
  ADD CONSTRAINT `bestellung_ibfk_2` FOREIGN KEY (`FK_ID_Kunde`) REFERENCES `Kunde` (`PK_ID_Kunde`),
  ADD CONSTRAINT `bestellung_ibfk_3` FOREIGN KEY (`FK_ID_Gericht`) REFERENCES `Gericht` (`PK_ID_Gericht`);

--
-- Constraints der Tabelle `Bewertung`
--
ALTER TABLE `Bewertung`
  ADD CONSTRAINT `bewertung_ibfk_1` FOREIGN KEY (`FK_ID_Kunde`) REFERENCES `Kunde` (`PK_ID_Kunde`) ON DELETE CASCADE,
  ADD CONSTRAINT `bewertung_ibfk_2` FOREIGN KEY (`FK_ID_Restaurant`) REFERENCES `Restaurant` (`PK_ID_Restaurant`) ON DELETE CASCADE;

--
-- Constraints der Tabelle `Gericht`
--
ALTER TABLE `Gericht`
  ADD CONSTRAINT `FK_Gericht_Kategorie` FOREIGN KEY (`FK_ID_Kategorie`) REFERENCES `Kategorie` (`PK_ID_Kategorie`),
  ADD CONSTRAINT `FK_ID_Restaurant` FOREIGN KEY (`FK_ID_Restaurant`) REFERENCES `Restaurant` (`PK_ID_Restaurant`);

--
-- Constraints der Tabelle `Gericht_Zutat`
--
ALTER TABLE `Gericht_Zutat`
  ADD CONSTRAINT `gericht_zutat_ibfk_1` FOREIGN KEY (`FK_ID_Gericht`) REFERENCES `Gericht` (`PK_ID_Gericht`) ON DELETE CASCADE,
  ADD CONSTRAINT `gericht_zutat_ibfk_2` FOREIGN KEY (`FK_ID_Zutat`) REFERENCES `Zutat` (`PK_ID_Zutat`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
