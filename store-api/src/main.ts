import { ValidationPipe } from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { EntityNotFouncExceptionFilter } from './exception-filters/entity-not-found-exception-filter';

async function bootstrap() {
   const app = await NestFactory.create(AppModule);
   app.useGlobalFilters(new EntityNotFouncExceptionFilter())
   app.useGlobalPipes(new ValidationPipe({ errorHttpStatusCode: 422 }))
   await app.listen(3000);
}
bootstrap();
