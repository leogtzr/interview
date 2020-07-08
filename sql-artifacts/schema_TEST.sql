-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema recruitment_interviews_test
-- -----------------------------------------------------
-- Recruitment interviews DB.
DROP SCHEMA IF EXISTS `recruitment_interviews_test` ;

-- -----------------------------------------------------
-- Schema recruitment_interviews_test
--
-- Recruitment interviews DB.
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `recruitment_interviews_test` ;
SHOW WARNINGS;
USE `recruitment_interviews_test` ;

-- -----------------------------------------------------
-- Table `recruitment_interviews_test`.`topic`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `recruitment_interviews_test`.`topic` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `recruitment_interviews_test`.`topic` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `topic` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `recruitment_interviews_test`.`level`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `recruitment_interviews_test`.`level` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `recruitment_interviews_test`.`level` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `recruitment_interviews_test`.`question`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `recruitment_interviews_test`.`question` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `recruitment_interviews_test`.`question` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `question` VARCHAR(1250) NOT NULL,
  `answer` VARCHAR(1000) NULL,
  `topic_id` INT NOT NULL,
  `level_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_question_topic_idx` (`topic_id` ASC),
  INDEX `fk_question_level1_idx` (`level_id` ASC),
  CONSTRAINT `fk_question_topic`
    FOREIGN KEY (`topic_id`)
    REFERENCES `recruitment_interviews_test`.`topic` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_question_level1`
    FOREIGN KEY (`level_id`)
    REFERENCES `recruitment_interviews_test`.`level` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `recruitment_interviews_test`.`candidate`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `recruitment_interviews_test`.`candidate` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `recruitment_interviews_test`.`candidate` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

SHOW WARNINGS;

-- -----------------------------------------------------
-- Table `recruitment_interviews_test`.`answer`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `recruitment_interviews_test`.`answer` ;

SHOW WARNINGS;
CREATE TABLE IF NOT EXISTS `recruitment_interviews_test`.`answer` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `result` INT NOT NULL,
  `comment` VARCHAR(1000) NULL,
  `question_id` INT NOT NULL,
  `candidate_id` INT NOT NULL,
  PRIMARY KEY (`id`, `candidate_id`),
  INDEX `fk_user_question_question1_idx` (`question_id` ASC),
  INDEX `fk_answer_candidate1_idx` (`candidate_id` ASC),
  CONSTRAINT `fk_user_question_question1`
    FOREIGN KEY (`question_id`)
    REFERENCES `recruitment_interviews_test`.`question` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_answer_candidate1`
    FOREIGN KEY (`candidate_id`)
    REFERENCES `recruitment_interviews_test`.`candidate` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;

SHOW WARNINGS;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

