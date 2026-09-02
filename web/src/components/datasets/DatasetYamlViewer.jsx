import React, { useState } from 'react';
import ClassSnapshotModal from './ClassSnapshotModal';
import { autoSuggestCategory } from './yolo_coco_classes';

export default function DatasetYamlViewer({ dataset, initialMappings = {}, onSaveMappings }) {
  const [activeModalClass, setActiveModalClass] = useState(null);
  const [classMappings, setClassMappings] = useState(initialMappings);
  const [saveStatus, setSaveStatus] = useState(null);

  React.useEffect(() => {
    setClassMappings(initialMappings);
  }, [dataset?.dataset_id, JSON.stringify(initialMappings)]);

  if (!dataset) {
    return (
      <div className="cyber-card">
        <div className="card-header"><span className="card-title">2. CLASS MAPPER</span></div>
        <p style={{ color: '#64748b', fontSize: '0.8rem' }}>No dataset selected.</p>
      </div>
    );
  }

  const handleSaveTarget = (origClass, targetClass) => {
    const updated = { ...classMappings, [origClass]: targetClass };
    setClassMappings(updated);
    if (onSaveMappings && dataset) onSaveMappings(dataset.dataset_id, updated);
    setSaveStatus('SAVED');
    setTimeout(() => setSaveStatus(null), 2500);
  };

  const rawClasses = dataset.classes || [];
  const targetNames = rawClasses
    .map(cls => classMappings[cls] || autoSuggestCategory(cls))
    .filter(c => c !== 'ignore' && c !== 'ignorar');
  const uniqueTargetClasses = Array.from(new Set(targetNames));
  const numClasses = uniqueTargetClasses.length > 0 ? uniqueTargetClasses.length : 1;

  const yamlString = `# Ultralytics YOLO26 Dataset Config
path: ${dataset.yaml_path ? dataset.yaml_path.replace('/data.yaml', '') : `/home/hades/datasets/${dataset.dataset_id}`}
train: train/images
val: valid/images
test: test/images

nc: ${numClasses}
names:
${uniqueTargetClasses.map((c, i) => `  ${i}: ${c}`).join('\n')}`;

  return (
    <div className="cyber-card" style={{ height: '380px', display: 'flex', flexDirection: 'column' }}>
      <div className="card-header" style={{ marginBottom: '0.4rem', paddingBottom: '0.35rem', flexShrink: 0 }}>
        <span className="card-title">2. CLASS MAPPER ({rawClasses.length} CLASSES)</span>
        <span className="badge-cyan">CLICK CLASS FOR SNAPSHOT</span>
      </div>

      <div style={{ flexShrink: 0, marginBottom: '0.3rem' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, minmax(0, 1fr))', gap: '0.3rem', maxHeight: '120px', overflowY: 'auto', paddingRight: '2px' }}>
          {rawClasses.map((cls) => {
            const target = classMappings[cls] || autoSuggestCategory(cls);
            return (
              <button key={cls} type="button" className="class-mapper-row" style={{ width: '100%', cursor: 'pointer', textAlign: 'left' }} onClick={() => setActiveModalClass(cls)}>
                <span className="class-mapper-label" title={cls}>🔍 "{cls}"</span>
                <span className="badge-yellow" style={{ fontSize: '0.62rem', padding: '0.1rem 0.35rem' }}>➔ {target}</span>
              </button>
            );
          })}
        </div>
      </div>

      <pre className="yaml-code-box" style={{ height: '90px', maxHeight: '90px', fontSize: '0.68rem', flexShrink: 0, marginTop: '0.3rem' }}>
        <code>{yamlString}</code>
      </pre>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '0.4rem', flexShrink: 0 }}>
        <span style={{ fontSize: '0.68rem', color: saveStatus ? 'var(--cb-green)' : 'var(--cb-cyan)', fontFamily: 'var(--font-mono)' }}>
          {saveStatus ? '✓ MAPPING SAVED' : `nc=${numClasses} (${uniqueTargetClasses.join(', ')})`}
        </span>
        <div style={{ display: 'flex', gap: '0.4rem' }}>
          <button className="cyber-action-btn secondary" style={{ padding: '0.35rem 0.8rem', fontSize: '0.7rem' }} onClick={() => navigator.clipboard.writeText(yamlString)}>COPY</button>
          <button className="cyber-action-btn" style={{ padding: '0.35rem 0.8rem', fontSize: '0.7rem' }} onClick={() => { if (onSaveMappings && dataset) onSaveMappings(dataset.dataset_id, classMappings); setSaveStatus('OK'); setTimeout(() => setSaveStatus(null), 2000); }}>APPLY</button>
        </div>
      </div>

      <ClassSnapshotModal
        isOpen={Boolean(activeModalClass)}
        onClose={() => setActiveModalClass(null)}
        datasetId={dataset.dataset_id}
        className={activeModalClass}
        currentTarget={classMappings[activeModalClass] || autoSuggestCategory(activeModalClass)}
        onSaveTarget={handleSaveTarget}
      />
    </div>
  );
}
