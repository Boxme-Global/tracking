ALTER TABLE "event" ADD COLUMN "otm_source" String;
ALTER TABLE "event" ADD COLUMN "otm_medium" String;
ALTER TABLE "event" ADD COLUMN "otm_campaign" String;
ALTER TABLE "event" ADD COLUMN "otm_position" String;


ALTER TABLE "page_view" ADD COLUMN "otm_source" String;
ALTER TABLE "page_view" ADD COLUMN "otm_medium" String;
ALTER TABLE "page_view" ADD COLUMN "otm_campaign" String;
ALTER TABLE "page_view" ADD COLUMN "otm_position" String;

ALTER TABLE "session" ADD COLUMN "otm_source" String;
ALTER TABLE "session" ADD COLUMN "otm_medium" String;
ALTER TABLE "session" ADD COLUMN "otm_campaign" String;
ALTER TABLE "session" ADD COLUMN "otm_position" String;