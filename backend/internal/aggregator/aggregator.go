
func (a *Aggregator) CleanUp() {
	slog.Info("Closing Kafka connection in aggregator")
	a.KafkaProducer.Close()
}
